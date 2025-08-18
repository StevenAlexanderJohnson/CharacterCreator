CREATE TABLE IF NOT EXISTS characters (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    owner_id INTEGER NOT NULl,
    name TEXT NOT NULL,
    bio TEXT NOT NULL DEFAULT '',
    background TEXT NOT NULL,
    class TEXT NOT NULL,
    level INTEGER NOT NULL,
    race_type TEXT NOT NULL,
    subrace_type TEXT, -- This column can be NULL
    race_move_speed INTEGER NOT NULL DEFAULT 0, -- This column can be NULL
    strength INTEGER NOT NULL,
    dexterity INTEGER NOT NULL,
    constitution INTEGER NOT NULL,
    intelligence INTEGER NOT NULL,
    wisdom INTEGER NOT NULL,
    charisma INTEGER NOT NULL,
    current_health_points INTEGER NOT NULL,
    
    FOREIGN KEY (owner_id) REFERENCES auth(id) ON DELETE CASCADE,
    FOREIGN KEY (class) REFERENCES classes(name),
    FOREIGN KEY (race_type) REFERENCES races(name)
);

CREATE TABLE IF NOT EXISTS character_proficiencies (
    character_id INTEGER NOT NULL,
    proficiency INTEGER NOT NULL,
    PRIMARY KEY (character_id, proficiency),
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE,
    FOREIGN KEY (proficiency) REFERENCES skills(name) ON DELETE CASCADE
);