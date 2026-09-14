package database

import (
	"context"
	"database/sql"
	"testing"
)

func TestOpenAppliesInitialMigrationIdempotently(t *testing.T) {
	db := openTestDB(t)

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&count); err != nil {
		t.Fatalf("query migrations: %v", err)
	}
	if count != 1 {
		t.Fatalf("migration count = %d, want 1", count)
	}
	for _, table := range []string{"resume_bullets", "skills", "bullet_skills", "tags", "bullet_tags"} {
		var exists bool
		if err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM sqlite_master WHERE type = 'table' AND name = ?)", table).Scan(&exists); err != nil {
			t.Fatalf("query table %s: %v", table, err)
		}
		if !exists {
			t.Errorf("expected metadata table %s to exist", table)
		}
	}

	if err := migrate(context.Background(), db); err != nil {
		t.Fatalf("second migrate() error = %v", err)
	}
}

func TestCandidateMetadataAssociationsAreIsolated(t *testing.T) {
	db := openTestDB(t)

	mustExec(t, db, `INSERT INTO candidate_profiles(id, display_name) VALUES ('candidate-a', 'A'), ('candidate-b', 'B')`)
	mustExec(t, db, `INSERT INTO experiences(id, candidate_id, organization, title) VALUES ('experience-a', 'candidate-a', 'Example', 'Engineer')`)
	mustExec(t, db, `INSERT INTO resume_bullets(id, candidate_id, experience_id, statement) VALUES ('bullet-a', 'candidate-a', 'experience-a', 'Improved reliability')`)
	mustExec(t, db, `INSERT INTO skills(id, candidate_id, name) VALUES ('skill-a', 'candidate-a', 'Go'), ('skill-b', 'candidate-b', 'React')`)
	mustExec(t, db, `INSERT INTO tags(id, candidate_id, name, kind) VALUES ('tag-a', 'candidate-a', 'Technical excellence', 'hat'), ('tag-b', 'candidate-b', 'Leadership', 'hat')`)

	mustExec(t, db, `INSERT INTO bullet_skills(candidate_id, bullet_id, skill_id, evidence) VALUES ('candidate-a', 'bullet-a', 'skill-a', 'Candidate supplied')`)
	mustExec(t, db, `INSERT INTO bullet_tags(candidate_id, bullet_id, tag_id, rationale) VALUES ('candidate-a', 'bullet-a', 'tag-a', 'Demonstrates engineering ownership')`)

	if _, err := db.Exec(`INSERT INTO bullet_skills(candidate_id, bullet_id, skill_id) VALUES ('candidate-a', 'bullet-a', 'skill-b')`); err == nil {
		t.Fatal("cross-candidate bullet skill insert succeeded, want foreign-key failure")
	}
	if _, err := db.Exec(`INSERT INTO bullet_tags(candidate_id, bullet_id, tag_id) VALUES ('candidate-b', 'bullet-a', 'tag-b')`); err == nil {
		t.Fatal("cross-candidate bullet tag insert succeeded, want foreign-key failure")
	}

	mustExec(t, db, `DELETE FROM candidate_profiles WHERE id = 'candidate-a'`)
	for _, table := range []string{"experiences", "resume_bullets", "bullet_skills", "bullet_tags"} {
		if got := rowCount(t, db, table); got != 0 {
			t.Errorf("%s row count after candidate deletion = %d, want 0", table, got)
		}
	}
	if got := rowCount(t, db, "skills"); got != 1 {
		t.Errorf("skills row count after candidate deletion = %d, want 1", got)
	}
	if got := rowCount(t, db, "tags"); got != 1 {
		t.Errorf("tags row count after candidate deletion = %d, want 1", got)
	}
}

func TestFinalResumeVersionsAreImmutable(t *testing.T) {
	db := openTestDB(t)
	seedResumeOwners(t, db)

	mustExec(t, db, `INSERT INTO resume_versions(
        id, candidate_id, job_target_id, source_snapshot_json, content_json,
        workflow_version, llm_provider, llm_model
    ) VALUES ('resume-a', 'candidate-a', 'job-a', '{"evidence_ids":["bullet-a"]}', '{"summary":"draft"}', 'v1', 'test', 'deterministic')`)

	mustExec(t, db, `UPDATE resume_versions SET content_json = '{"summary":"final"}', status = 'final' WHERE id = 'resume-a'`)

	if _, err := db.Exec(`UPDATE resume_versions SET content_json = '{"summary":"changed"}' WHERE id = 'resume-a'`); err == nil {
		t.Fatal("final resume update succeeded, want immutable-version failure")
	}
	if _, err := db.Exec(`DELETE FROM resume_versions WHERE id = 'resume-a'`); err == nil {
		t.Fatal("final resume deletion succeeded, want immutable-version failure")
	}

	var content string
	if err := db.QueryRow(`SELECT content_json FROM resume_versions WHERE id = 'resume-a'`).Scan(&content); err != nil {
		t.Fatalf("query final resume: %v", err)
	}
	if content != `{"summary":"final"}` {
		t.Errorf("final resume content = %q, want finalized content", content)
	}
}

func TestApplicationsRequireFinalResumeVersion(t *testing.T) {
	db := openTestDB(t)
	seedResumeOwners(t, db)

	mustExec(t, db, `INSERT INTO resume_versions(
        id, candidate_id, job_target_id, source_snapshot_json, content_json,
        workflow_version, llm_provider, llm_model
    ) VALUES
        ('resume-draft', 'candidate-a', 'job-a', '{}', '{}', 'v1', 'test', 'deterministic'),
        ('resume-final', 'candidate-a', 'job-a', '{}', '{}', 'v1', 'test', 'deterministic')`)
	mustExec(t, db, `UPDATE resume_versions SET status = 'final' WHERE id = 'resume-final'`)

	if _, err := db.Exec(`INSERT INTO applications(id, job_target_id, resume_version_id) VALUES ('application-draft', 'job-a', 'resume-draft')`); err == nil {
		t.Fatal("application linked to draft resume, want finalized-version failure")
	}
	mustExec(t, db, `INSERT INTO applications(id, job_target_id, resume_version_id) VALUES ('application-final', 'job-a', 'resume-final')`)

	if _, err := db.Exec(`UPDATE applications SET resume_version_id = 'resume-draft' WHERE id = 'application-final'`); err == nil {
		t.Fatal("application changed to draft resume, want finalized-version failure")
	}
}

func TestResumeRevisionLineageUsesFinalParentFromSameCandidate(t *testing.T) {
	db := openTestDB(t)
	seedResumeOwners(t, db)
	mustExec(t, db, `INSERT INTO candidate_profiles(id, display_name) VALUES ('candidate-b', 'B')`)

	mustExec(t, db, `INSERT INTO resume_versions(
        id, candidate_id, job_target_id, source_snapshot_json, content_json,
        workflow_version, llm_provider, llm_model
    ) VALUES
        ('resume-parent', 'candidate-a', 'job-a', '{}', '{}', 'v1', 'test', 'deterministic'),
        ('resume-other-candidate', 'candidate-b', 'job-a', '{}', '{}', 'v1', 'test', 'deterministic')`)

	if _, err := db.Exec(`INSERT INTO resume_versions(
        id, candidate_id, job_target_id, parent_version_id, source_snapshot_json, content_json,
        workflow_version, llm_provider, llm_model
    ) VALUES ('child-of-draft', 'candidate-a', 'job-a', 'resume-parent', '{}', '{}', 'v1', 'test', 'deterministic')`); err == nil {
		t.Fatal("revision linked to draft parent, want finalized-parent failure")
	}

	mustExec(t, db, `UPDATE resume_versions SET status = 'final' WHERE id IN ('resume-parent', 'resume-other-candidate')`)
	mustExec(t, db, `INSERT INTO resume_versions(
        id, candidate_id, job_target_id, parent_version_id, source_snapshot_json, content_json,
        workflow_version, llm_provider, llm_model
    ) VALUES ('resume-child', 'candidate-a', 'job-a', 'resume-parent', '{}', '{}', 'v2', 'test', 'deterministic')`)

	if _, err := db.Exec(`INSERT INTO resume_versions(
        id, candidate_id, job_target_id, parent_version_id, source_snapshot_json, content_json,
        workflow_version, llm_provider, llm_model
    ) VALUES ('cross-candidate-child', 'candidate-a', 'job-a', 'resume-other-candidate', '{}', '{}', 'v2', 'test', 'deterministic')`); err == nil {
		t.Fatal("revision linked across candidates, want ownership failure")
	}
	if _, err := db.Exec(`DELETE FROM resume_versions WHERE id = 'resume-parent'`); err == nil {
		t.Fatal("finalized lineage parent deletion succeeded, want failure")
	}
	if _, err := db.Exec(`DELETE FROM job_targets WHERE id = 'job-a'`); err == nil {
		t.Fatal("job target deletion succeeded while resume snapshots reference it, want restricted deletion")
	}
}

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	dsn := "file:" + t.TempDir() + "/resonance.db?_pragma=foreign_keys(1)"
	db, err := Open(context.Background(), dsn)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("close database: %v", err)
		}
	})
	return db
}

func seedResumeOwners(t *testing.T, db *sql.DB) {
	t.Helper()
	mustExec(t, db, `INSERT INTO candidate_profiles(id, display_name) VALUES ('candidate-a', 'A')`)
	mustExec(t, db, `INSERT INTO job_targets(id, company, role_title, description_snapshot) VALUES ('job-a', 'Example', 'Engineer', 'Build reliable software')`)
}

func mustExec(t *testing.T, db *sql.DB, query string) {
	t.Helper()
	if _, err := db.Exec(query); err != nil {
		t.Fatalf("exec %q: %v", query, err)
	}
}

func rowCount(t *testing.T, db *sql.DB, table string) int {
	t.Helper()
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count); err != nil {
		t.Fatalf("count %s: %v", table, err)
	}
	return count
}
