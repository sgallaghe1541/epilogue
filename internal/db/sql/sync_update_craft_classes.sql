INSERT INTO craft_classes (template, class, description, active)
SELECT template, class, description, 1
FROM temp_craft_classes
WHERE template NOT IN (
    SELECT template
    FROM craft_classes
);

UPDATE craft_classes
SET active = 0
WHERE template NOT IN (
    SELECT template
    FROM temp_craft_classes
);