DROP TABLE IF EXISTS temp_craftTemplates;
DROP TABLE IF EXISTS temp_jobs;
DROP TABLE IF EXISTS temp_phases;
DROP TABLE IF EXISTS temp_employees;
DROP TABLE IF EXISTS temp_equipment;

CREATE TABLE IF NOT EXISTS craftClasses (
    template INTEGER,
    class TEXT,
    description TEXT,
    active INTEGER, -- 0 is inactive, 1 is active
    PRIMARY KEY(template, class)
);

CREATE TEMPORARY TABLE temp_craftClasses (
    template INTEGER,
    class TEXT,
    description TEXT,
    PRIMARY KEY(template, class)
);

CREATE TABLE IF NOT EXISTS jobs (
    job TEXT NOT NULL PRIMARY KEY,
    description TEXT,
    state TEXT,
    certified TEXT,
    template INTEGER,
    active INTEGER, -- 0 is inactive, 1 is active
    FOREIGN KEY(template) REFERENCES craftClasses(template)
);

CREATE TEMPORARY TABLE temp_jobs (
    job TEXT NOT NULL PRIMARY KEY,
    description TEXT,
    state TEXT,
    certified TEXT,
    template INTEGER,
    FOREIGN KEY(template) REFERENCES temp_craftClasses(template)
);

CREATE TABLE IF NOT EXISTS phases (
    job TEXT NOT NULL,
    phase TEXT NOT NULL,
    description TEXT,
    active INTEGER, -- 0 is inactive, 1 is active
    PRIMARY KEY(job, phase),
    FOREIGN KEY(job) REFERENCES jobs(job)
);

CREATE TEMPORARY TABLE temp_phases (
    job TEXT NOT NULL,
    phase TEXT NOT NULL,
    description TEXT,
    PRIMARY KEY(job, phase),
    FOREIGN KEY(job) REFERENCES temp_jobs(job)
);

CREATE TABLE IF NOT EXISTS employees (
    employee TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    class TEXT,
    salaried INTEGER, -- 0 is hourly, 1 is salaried
    department TEXT,
    active INTEGER -- 0 is inactive, 1 is active
);

CREATE TEMPORARY TABLE temp_employees (
    employee TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    class TEXT,
    salaried INTEGER, -- 0 is hourly, 1 is salaried
    department TEXT
);

CREATE TABLE IF NOT EXISTS equipment (
    equipment TEXT PRIMARY KEY,
    description TEXT,
    department TEXT,
    active INTEGER -- 0 is inactive, 1 is active
);

CREATE TEMPORARY TABLE temp_equipment (
    equipment TEXT PRIMARY KEY,
    description TEXT,
    department TEXT
);
