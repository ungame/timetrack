CREATE TABLE IF NOT EXISTS categories (
    id INT AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT categories_id_pk PRIMARY KEY(id)
)
ENGINE = INNODB
DEFAULT CHARSET = UTF8;

INSERT INTO categories (name, description) VALUES
('meeting', 'daily, 1:1, refinement, retro'),
('coding', 'features, bugs, tests'),
('review', 'pull requests'),
('pm', 'private messages'),
('task', '')
;

CREATE TABLE IF NOT EXISTS activities (
    id BIGINT AUTO_INCREMENT,
    category_id INT NOT NULL,
    description TEXT NOT NULL,
    status CHAR(1) DEFAULT '1',
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    finished_at TIMESTAMP NULL,
    CONSTRAINT activities_id_pk PRIMARY KEY(id),
    CONSTRAINT activities_category_id_fk
    FOREIGN KEY (category_id)
    REFERENCES categories(id)
)
ENGINE = INNODB
DEFAULT CHARSET = UTF8;