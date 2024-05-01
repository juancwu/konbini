CREATE TABLE email_verifications (
    id SERIAL NOT NULL PRIMARY KEY,
    ref_id VARCHAR(16) NOT NULL UNIQUE,
    status STATUS NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);