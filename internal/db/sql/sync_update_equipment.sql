INSERT INTO equipment (equipment, description, department, active)
SELECT equipment, description, department, 1
FROM temp_equipment
WHERE equipment NOT IN (
    SELECT equipment
    FROM equipment
);

UPDATE equipment
SET active = 1
WHERE active = 0
AND equipment IN (
    SELECT equipment
    FROM temp_equipment
);

UPDATE equipment
SET active = 0
WHERE equipment NOT IN (
    SELECT equipment
    FROM temp_equipment
);

UPDATE equipment
SET
    description = temp_equipment.description,
    department = temp_equipment.department
FROM (
    SELECT equipment, description, department
    FROM temp_equipment 
) AS temp_equipment
WHERE temp_equipment.equipment = equipment.equipment;