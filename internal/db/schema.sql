CREATE TABLE IF NOT EXISTS divisions (
    divisionid INT PRIMARY KEY,
    descr TEXT
);

CREATE TABLE IF NOT EXISTS permissionlevels (
    permissionid INT PRIMARY KEY,
    descr TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    email TEXT,
    username TEXT,
    division INTEGER,
    permissionlevel INTEGER,
    FOREIGN KEY(division) REFERENCES divisions(divisionid),
    FOREIGN KEY(permissionlevel) REFERENCES permissionlevels(permissionid)
);

CREATE TABLE IF NOT EXISTS reports (
    reportid INTEGER PRIMARY KEY ASC,
    reportname TEXT,
    reporturl TEXT,
    permission INTEGER,
    FOREIGN KEY(permission) REFERENCES permissionlevel(permissionid)
);

CREATE TABLE IF NOT EXISTS parametertypes (
    paramtype TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS parameters (
    paramname TEXT NOT NULL,
    paramdesc TEXT NOT NULL,
    reportid INTEGER,
    paramtype TEXT,
    FOREIGN KEY(reportid) REFERENCES reports(reportid),
    FOREIGN KEY(paramtype) REFERENCES parametertypes(paramtype),
    PRIMARY KEY(paramname, reportid)
);

CREATE TABLE IF NOT EXISTS icons (
    iconid INTEGER PRIMARY KEY ASC,
    iconurl TEXT,
    tooltip TEXT,
    permission INTEGER,
    FOREIGN KEY(permission) REFERENCES permissionlevel(permissionid)
);

CREATE TABLE IF NOT EXISTS iconpaths (
    iconpath TEXT,
    iconid INTEGER,
    FOREIGN KEY(iconid) REFERENCES icons(iconid)
);