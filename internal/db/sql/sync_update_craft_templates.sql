INSERT INTO craftTemplates (template, class, description, active)
SELECT template, class, description, 1
FROM temp_craftTemplates
WHERE template NOT IN (
    SELECT template
    FROM craftTemplates
);

UPDATE craftTemplates
SET active = 0
WHERE template NOT IN (
    SELECT template
    FROM temp_craftTemplates
);