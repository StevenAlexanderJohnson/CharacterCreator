CREATE TABLE IF NOT EXISTS skills (
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS classes (
    name TEXT NOT NULL UNIQUE
);

INSERT INTO skills (name) VALUES
('Acrobatics'),
('Animal Handling'),
('Arcana'),
('Athletics'),
('Deception'),
('History'),
('Insight'),
('Intimidation'),
('Investigation'),
('Medicine'),
('Nature'),
('Perception'),
('Performance'),
('Persuasion'),
('Religion'),
('Sleight of Hand'),
('Stealth'),
('Survival');

INSERT INTO classes (name) VALUES
('Barbarian'),
('Bard'),
('Cleric'),
('Druid'),
('Fighter'),
('Monk'),
('Paladin'),
('Ranger'),
('Rogue'),
('Sorcerer'),
('Warlock'),
('Wizard'),
('Commoner');