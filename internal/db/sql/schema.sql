CREATE TABLE IF NOT EXISTS departments (
    departmentid TEXT PRIMARY KEY,
    description TEXT KEY
);

CREATE TABLE IF NOT EXISTS roles (
    roleid INTEGER PRIMARY KEY,
    role TEXT
);

CREATE TABLE IF NOT EXISTS users (
    userid INTEGER PRIMARY KEY ASC,
    email TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    active INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS user_permissions (
    userid INTEGER,
    department TEXT,
    roleid INTEGER,
    FOREIGN KEY(userid) REFERENCES users(userid),
    FOREIGN KEY(department) REFERENCES departments(departmentid),
    FOREIGN KEY(roleid) REFERENCES roles(roleid),
    PRIMARY KEY(userid, department, roleid)
);

CREATE TABLE IF NOT EXISTS resource_types (
    resourcetype TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS resources (
    resourceid INTEGER PRIMARY KEY ASC,
    resourcetype TEXT,
    resourcename TEXT,
    resourceurl TEXT,
    FOREIGN KEY(resourcetype) REFERENCES resource_types(resourcetype)
);

CREATE TABLE IF NOT EXISTS resource_permissions (
    resourceid INTEGER,
    department TEXT,
    roleid INTEGER,
    FOREIGN KEY(resourceid) REFERENCES resources(resourceid),
    FOREIGN KEY(department) REFERENCES departments(departmentid),
    FOREIGN KEY(roleid) REFERENCES roles(roleid),
    PRIMARY KEY(resourceid, department, roleid)
);

CREATE TABLE IF NOT EXISTS parameter_types (
    paramtype TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS parameters (
    paramid INTEGER ASC,
    paramname TEXT NOT NULL,
    paramdesc TEXT NOT NULL,
    resourceid INTEGER,
    paramtype TEXT,
    FOREIGN KEY(resourceid) REFERENCES resources(resourceid),
    FOREIGN KEY(paramtype) REFERENCES parameter_types(paramtype),
    PRIMARY KEY(paramname, resourceid)
);

CREATE TABLE IF NOT EXISTS param_select_list (
    paramid INTEGER,
    selectdisplayas TEXT,
    selectvalue TEXT,
    FOREIGN KEY(paramid) REFERENCES parameters(paramid)
    PRIMARY KEY(paramid, selectdisplayas)
);

CREATE TABLE IF NOT EXISTS icons (
    iconid INTEGER PRIMARY KEY ASC,
    iconurl TEXT,
    tooltip TEXT,
    roleid INTEGER,
    FOREIGN KEY(roleid) REFERENCES roles(roleid)
);

CREATE TABLE IF NOT EXISTS iconpaths (
    iconpath TEXT,
    iconid INTEGER,
    FOREIGN KEY(iconid) REFERENCES icons(iconid)
);

CREATE TABLE IF NOT EXISTS sessions (
	token TEXT PRIMARY KEY,
	data BLOB NOT NULL,
	expiry REAL NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions(expiry);

CREATE TABLE IF NOT EXISTS refreshtokens (
    userid INT,
    token TEXT,
    expires REAL,
    FOREIGN KEY(userid) REFERENCES users(userid)
);

CREATE TABLE IF NOT EXISTS craft_classes (
    template INTEGER,
    class TEXT,
    description TEXT,
    active INTEGER, -- 0 is inactive, 1 is active
    PRIMARY KEY(template, class)
);

CREATE TABLE IF NOT EXISTS jobs (
    job TEXT NOT NULL PRIMARY KEY,
    description TEXT,
    state TEXT,
    certified TEXT,
    template INTEGER,
    department TEXT,
    active INTEGER, -- 0 is inactive, 1 is active
    FOREIGN KEY(template) REFERENCES craft_classes(template)
);

CREATE TABLE IF NOT EXISTS phases (
    job TEXT NOT NULL,
    phase TEXT NOT NULL,
    description TEXT,
    active INTEGER, -- 0 is inactive, 1 is active
    PRIMARY KEY(job, phase),
    FOREIGN KEY(job) REFERENCES jobs(job)
);

CREATE TABLE IF NOT EXISTS employees (
    employee TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    class TEXT,
    salaried INTEGER, -- 0 is hourly, 1 is salaried
    department TEXT,
    active INTEGER -- 0 is inactive, 1 is active
);

CREATE TABLE IF NOT EXISTS equipment (
    equipment TEXT PRIMARY KEY,
    description TEXT,
    department TEXT,
    active INTEGER -- 0 is inactive, 1 is active
);

CREATE TABLE IF NOT EXISTS time_card_status (
    tcstatus TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS time_card_headers (
    id INTEGER PRIMARY KEY ASC,
    job TEXT,
    workdate DATE,
    createdby INTEGER,
    tcstatus TEXT,
    lastmodified DATE,
    modifiedby INTEGER,
    FOREIGN KEY(job) REFERENCES jobs(job),
    FOREIGN KEY(createdby) REFERENCES users(userid),
    FOREIGN KEY(tcstatus) REFERENCES time_card_status(tcstatus),
    FOREIGN KEY(modifiedby) REFERENCES users(userid)
);

CREATE TABLE IF NOT EXISTS time_card_employees (
    tceid INTEGER PRIMARY KEY ASC,
    tchid INTEGER,
    employee TEXT,
    fullname TEXT,
    workdate DATE,
    job TEXT,
    phase TEXT,
    class TEXT,
    paycode TEXT,
    tcehours NUMERIC,
    FOREIGN KEY(tchid) REFERENCES time_card_headers(id),
    FOREIGN KEY(job) REFERENCES jobs(job),
    FOREIGN KEY(employee) REFERENCES employees(employee)
);

CREATE TABLE IF NOT EXISTS earn_codes (
    description TEXT PRIMARY KEY
);