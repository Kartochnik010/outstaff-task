CREATE TABLE IF NOT EXISTS music (
    id BIGSERIAL PRIMARY KEY,
	group_name       text NOT NULL,
	song        text NOT NULL,
	link        text NOT NULL,
	text        text NOT NULL,
	release_date TIMESTAMP NOT NULL
);
