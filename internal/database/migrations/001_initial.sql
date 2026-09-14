CREATE TABLE IF NOT EXISTS candidate_profiles (
    id TEXT PRIMARY KEY,
    display_name TEXT NOT NULL DEFAULT '',
    headline TEXT NOT NULL DEFAULT '',
    summary TEXT NOT NULL DEFAULT '',
    contact_json TEXT NOT NULL DEFAULT '{}',
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

CREATE TABLE IF NOT EXISTS experiences (
    id TEXT PRIMARY KEY,
    candidate_id TEXT NOT NULL REFERENCES candidate_profiles(id) ON DELETE CASCADE,
    organization TEXT NOT NULL,
    title TEXT NOT NULL,
    location TEXT NOT NULL DEFAULT '',
    started_on TEXT,
    ended_on TEXT,
    is_current INTEGER NOT NULL DEFAULT 0 CHECK (is_current IN (0, 1)),
    source_notes TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

CREATE TABLE IF NOT EXISTS resume_bullets (
    id TEXT PRIMARY KEY,
    experience_id TEXT NOT NULL REFERENCES experiences(id) ON DELETE CASCADE,
    statement TEXT NOT NULL,
    situation TEXT NOT NULL DEFAULT '',
    action TEXT NOT NULL DEFAULT '',
    result TEXT NOT NULL DEFAULT '',
    evidence TEXT NOT NULL DEFAULT '',
    sort_order INTEGER NOT NULL DEFAULT 0,
    is_archived INTEGER NOT NULL DEFAULT 0 CHECK (is_archived IN (0, 1)),
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

CREATE TABLE IF NOT EXISTS skills (
    id TEXT PRIMARY KEY,
    candidate_id TEXT NOT NULL REFERENCES candidate_profiles(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    category TEXT NOT NULL DEFAULT '',
    evidence TEXT NOT NULL DEFAULT '',
    source TEXT NOT NULL DEFAULT 'candidate' CHECK (source IN ('candidate', 'imported', 'confirmed-inference')),
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    UNIQUE (candidate_id, name)
);

CREATE TABLE IF NOT EXISTS bullet_skills (
    bullet_id TEXT NOT NULL REFERENCES resume_bullets(id) ON DELETE CASCADE,
    skill_id TEXT NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    source TEXT NOT NULL DEFAULT 'candidate' CHECK (source IN ('candidate', 'imported', 'confirmed-inference')),
    evidence TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (bullet_id, skill_id)
);

CREATE TABLE IF NOT EXISTS tags (
    id TEXT PRIMARY KEY,
    candidate_id TEXT NOT NULL REFERENCES candidate_profiles(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    kind TEXT NOT NULL CHECK (kind IN ('hat', 'strength', 'domain', 'outcome', 'keyword')),
    description TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    UNIQUE (candidate_id, kind, name)
);

CREATE TABLE IF NOT EXISTS bullet_tags (
    bullet_id TEXT NOT NULL REFERENCES resume_bullets(id) ON DELETE CASCADE,
    tag_id TEXT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    source TEXT NOT NULL DEFAULT 'candidate' CHECK (source IN ('candidate', 'imported', 'confirmed-inference')),
    rationale TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (bullet_id, tag_id)
);

CREATE TABLE IF NOT EXISTS job_targets (
    id TEXT PRIMARY KEY,
    company TEXT NOT NULL,
    role_title TEXT NOT NULL,
    source_url TEXT NOT NULL DEFAULT '',
    description_snapshot TEXT NOT NULL,
    captured_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

CREATE TABLE IF NOT EXISTS resume_versions (
    id TEXT PRIMARY KEY,
    candidate_id TEXT NOT NULL REFERENCES candidate_profiles(id) ON DELETE RESTRICT,
    job_target_id TEXT REFERENCES job_targets(id) ON DELETE SET NULL,
    parent_version_id TEXT REFERENCES resume_versions(id) ON DELETE SET NULL,
    status TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'final')),
    source_snapshot_json TEXT NOT NULL,
    gap_analysis_json TEXT NOT NULL DEFAULT '{}',
    content_json TEXT NOT NULL,
    rendered_path TEXT NOT NULL DEFAULT '',
    workflow_version TEXT NOT NULL,
    llm_provider TEXT NOT NULL,
    llm_model TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

CREATE TABLE IF NOT EXISTS applications (
    id TEXT PRIMARY KEY,
    job_target_id TEXT NOT NULL REFERENCES job_targets(id) ON DELETE RESTRICT,
    resume_version_id TEXT NOT NULL REFERENCES resume_versions(id) ON DELETE RESTRICT,
    status TEXT NOT NULL DEFAULT 'planned' CHECK (
        status IN ('planned', 'applied', 'screening', 'interviewing', 'offer', 'rejected', 'withdrawn', 'closed')
    ),
    applied_at TEXT,
    notes TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

CREATE INDEX IF NOT EXISTS idx_experiences_candidate ON experiences(candidate_id);
CREATE INDEX IF NOT EXISTS idx_resume_bullets_experience ON resume_bullets(experience_id, sort_order);
CREATE INDEX IF NOT EXISTS idx_skills_candidate ON skills(candidate_id);
CREATE INDEX IF NOT EXISTS idx_tags_candidate_kind ON tags(candidate_id, kind, name);
CREATE INDEX IF NOT EXISTS idx_bullet_skills_skill ON bullet_skills(skill_id, bullet_id);
CREATE INDEX IF NOT EXISTS idx_bullet_tags_tag ON bullet_tags(tag_id, bullet_id);
CREATE INDEX IF NOT EXISTS idx_resume_versions_candidate ON resume_versions(candidate_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_applications_status ON applications(status, updated_at DESC);
