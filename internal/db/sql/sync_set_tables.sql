DROP TABLE IF EXISTS temp_craft_classes;
DROP TABLE IF EXISTS temp_jobs;
DROP TABLE IF EXISTS temp_phases;
DROP TABLE IF EXISTS temp_employees;
DROP TABLE IF EXISTS temp_equipment;

CREATE TEMPORARY TABLE temp_craft_classes (
    template INTEGER,
    class TEXT,
    description TEXT,
    PRIMARY KEY(template, class)
);

CREATE TEMPORARY TABLE temp_jobs (
    job TEXT NOT NULL PRIMARY KEY,
    description TEXT,
    state TEXT,
    certified TEXT,
    template INTEGER,
    department TEXT,
    FOREIGN KEY(template) REFERENCES temp_craft_classes(template)
);

CREATE TEMPORARY TABLE temp_phases (
    job TEXT NOT NULL,
    phase TEXT NOT NULL,
    description TEXT,
    PRIMARY KEY(job, phase),
    FOREIGN KEY(job) REFERENCES temp_jobs(job)
);

CREATE TEMPORARY TABLE temp_employees (
    employee TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    class TEXT,
    salaried INTEGER, -- 0 is hourly, 1 is salaried
    department TEXT
);

CREATE TEMPORARY TABLE temp_equipment (
    equipment TEXT PRIMARY KEY,
    description TEXT,
    department TEXT
);
