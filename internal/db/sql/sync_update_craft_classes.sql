INSERT INTO craftClasses (template, class, description, active)
SELECT template, class, description, 1
FROM temp_craftClasses
WHERE template NOT IN (
    SELECT template
    FROM craftClasses
);

UPDATE craftClasses
SET active = 0
WHERE template NOT IN (
    SELECT template
    FROM temp_craftClasses
);