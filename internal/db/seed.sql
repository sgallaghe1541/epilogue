INSERT INTO permissionlevels (permissionid, descr)
VALUES 
    (1, 'Any'),
    (2, 'No Payroll'),
    (3, 'Division Head'),
    (5, 'Admin');

INSERT INTO reports (reportid, reportname, reporturl, permission)
VALUES (1, 'All Job Hours', 'alljobhours', 2);

INSERT INTO parametertypes (paramtype)
VALUES 
    ('date');

INSERT INTO parameters (paramname, paramdesc, reportid, paramtype)
VALUES
    ('startwedate', 'Beginning WE Date', 1, 'date'),
    ('endwedate', 'Ending WE Date', 1, 'date');