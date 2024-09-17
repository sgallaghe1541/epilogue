INSERT INTO jobs (job, description, state, certified, template, active)
SELECT job, description, state, certified, template, 1
FROM temp_jobs
WHERE job NOT IN (
    SELECT job
    FROM jobs
);

UPDATE jobs
SET active = 1
WHERE active = 0
AND job IN (
    SELECT job
    FROM temp_jobs 
);

UPDATE jobs
SET active = 0
WHERE job NOT IN (
    SELECT job
    FROM temp_jobs 
);

DROP TABLE IF EXISTS temp_jobs;