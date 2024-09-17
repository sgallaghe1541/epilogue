INSERT INTO phases (job, phase, description, active)
SELECT job, phase, description, 1
FROM temp_phases
WHERE job NOT IN (
    SELECT job
    FROM phases
);

UPDATE phases
SET active = 1
WHERE active = 0
AND job NOT IN (
    SELECT job
    FROM temp_phases 
);

UPDATE phases
SET active = 0
WHERE job NOT IN (
    SELECT job
    FROM temp_phases 
);

DROP TABLE IF EXISTS temp_phases;