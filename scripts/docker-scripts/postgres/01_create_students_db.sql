DO
$$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'students') THEN
      RAISE NOTICE 'Creating database students...'; -- Логирование
      CREATE DATABASE students;
   ELSE
      RAISE NOTICE 'Database students already exists'; -- Логирование
   END IF;
END
$$
