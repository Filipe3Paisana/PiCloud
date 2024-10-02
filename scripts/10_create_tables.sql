-- Tabela 'Users'
CREATE TABLE IF NOT EXISTS Users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Tabela 'Files'
CREATE TABLE IF NOT EXISTS Files (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    size INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    user_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE
);

-- Tabela 'FileFragments'
CREATE TABLE IF NOT EXISTS FileFragments (
    fragment_id SERIAL PRIMARY KEY,
    file_id INT NOT NULL,
    fragment_hashes VARCHAR(255) NOT NULL,
    fragment_order INT NOT NULL,
    FOREIGN KEY (file_id) REFERENCES Files(id) ON DELETE CASCADE
);

-- Tabela 'Nodes'
CREATE TABLE IF NOT EXISTS Nodes (
    id SERIAL PRIMARY KEY,
    node_address VARCHAR(255) NOT NULL,
    location VARCHAR(255) NOT NULL,
    capacity INT NOT NULL,
    available_capacity INT NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Tabela 'FragmentLocation'
CREATE TABLE IF NOT EXISTS FragmentLocation (
    fragment_id INT NOT NULL,
    node_id INT NOT NULL,
    stored_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (fragment_id, node_id),
    FOREIGN KEY (fragment_id) REFERENCES FileFragments(fragment_id) ON DELETE CASCADE,
    FOREIGN KEY (node_id) REFERENCES Nodes(id) ON DELETE CASCADE
);
