INSERT INTO employees (employee, name, class, salaried, department, active)
SELECT employee, name, class, salaried, department, 1
FROM temp_employees
WHERE employee NOT IN (
    SELECT employee
    FROM employees
);

UPDATE employees
SET active = 1
WHERE active = 0
AND employee IN (
    SELECT employee
    FROM temp_phases 
);

UPDATE employees
SET active = 0
WHERE employee NOT IN (
    SELECT employee
    FROM temp_employees
);

UPDATE employees
SET
    class = temp_employees.class,
    salaried = temp_employees.salaried,
    department = temp_employees.department
FROM (
    SELECT employee, class, salaried, department
    FROM temp_employees
) AS temp_employees
WHERE temp_employees.employee = employees.employee;