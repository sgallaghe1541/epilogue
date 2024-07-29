INSERT INTO divisions (divisionid, descr)
VALUES
    ('01', 'Grading'),
    ('02', 'Paving'),
    ('06', 'Bridge');

INSERT INTO permissionlevels (permissionid, descr)
VALUES 
    (1, 'Any'),
    (2, 'No Payroll'),
    (3, 'Division Head'),
    (4, 'Payroll'),
    (5, 'Admin');

INSERT INTO reports (reportid, reportname, reporturl, permission)
VALUES 
    (1, 'All Job Hours', 'alljobhours', 2),
    (2, 'Employees for Fringe', 'employeesforfringe', 4);

INSERT INTO parametertypes (paramtype)
VALUES 
    ('date'),
    ('division');

INSERT INTO parameters (paramname, paramdesc, reportid, paramtype)
VALUES
    ('startwedate', 'Beginning WE Date', 1, 'date'),
    ('endwedate', 'Ending WE Date', 1, 'date'),
    ('division', 'Division', 1, 'division'),
    ('division', 'Division', 2, 'division');

INSERT INTO icons (iconid, iconurl, tooltip, permission)
VALUES
    (1, '/reports/', 'Reports', 1),
    (2, '/timeentry/', 'Time Entry', 2),
    (3, '/tools/', 'Tools', 4),
    (4, '/settings/', 'Settings', 1),
    (5, '/admin/', 'Admin', 5);

INSERT INTO iconpaths (iconpath, iconid)
VALUES
    ('M19.5 14.25v-2.625a3.375 3.375 0 0 0-3.375-3.375h-1.5A1.125 1.125 0 0 1 13.5 7.125v-1.5a3.375 3.375 0 0 0-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 0 0-9-9Z', 1),
    ('M12 6v6h4.5m4.5 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z', 2),
    ('M11.42 15.17 17.25 21A2.652 2.652 0 0 0 21 17.25l-5.877-5.877M11.42 15.17l2.496-3.03c.317-.384.74-.626 1.208-.766M11.42 15.17l-4.655 5.653a2.548 2.548 0 1 1-3.586-3.586l6.837-5.63m5.108-.233c.55-.164 1.163-.188 1.743-.14a4.5 4.5 0 0 0 4.486-6.336l-3.276 3.277a3.004 3.004 0 0 1-2.25-2.25l3.276-3.276a4.5 4.5 0 0 0-6.336 4.486c.091 1.076-.071 2.264-.904 2.95l-.102.085m-1.745 1.437L5.909 7.5H4.5L2.25 3.75l1.5-1.5L7.5 4.5v1.409l4.26 4.26m-1.745 1.437 1.745-1.437m6.615 8.206L15.75 15.75M4.867 19.125h.008v.008h-.008v-.008Z', 3),
    ('M9.594 3.94c.09-.542.56-.94 1.11-.94h2.593c.55 0 1.02.398 1.11.94l.213 1.281c.063.374.313.686.645.87.074.04.147.083.22.127.325.196.72.257 1.075.124l1.217-.456a1.125 1.125 0 0 1 1.37.49l1.296 2.247a1.125 1.125 0 0 1-.26 1.431l-1.003.827c-.293.241-.438.613-.43.992a7.723 7.723 0 0 1 0 .255c-.008.378.137.75.43.991l1.004.827c.424.35.534.955.26 1.43l-1.298 2.247a1.125 1.125 0 0 1-1.369.491l-1.217-.456c-.355-.133-.75-.072-1.076.124a6.47 6.47 0 0 1-.22.128c-.331.183-.581.495-.644.869l-.213 1.281c-.09.543-.56.94-1.11.94h-2.594c-.55 0-1.019-.398-1.11-.94l-.213-1.281c-.062-.374-.312-.686-.644-.87a6.52 6.52 0 0 1-.22-.127c-.325-.196-.72-.257-1.076-.124l-1.217.456a1.125 1.125 0 0 1-1.369-.49l-1.297-2.247a1.125 1.125 0 0 1 .26-1.431l1.004-.827c.292-.24.437-.613.43-.991a6.932 6.932 0 0 1 0-.255c.007-.38-.138-.751-.43-.992l-1.004-.827a1.125 1.125 0 0 1-.26-1.43l1.297-2.247a1.125 1.125 0 0 1 1.37-.491l1.216.456c.356.133.751.072 1.076-.124.072-.044.146-.086.22-.128.332-.183.582-.495.644-.869l.214-1.28Z', 4),
    ('M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z', 4),
    ('M17.25 6.75 22.5 12l-5.25 5.25m-10.5 0L1.5 12l5.25-5.25m7.5-3-4.5 16.5', 5);

INSERT INTO users (userid, email, name, division, permissionlevel)
VALUES (1, "sgallagher@talleyconstruction.net", "Seth Gallagher", "00", 5);