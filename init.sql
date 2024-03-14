CREATE TABLE movie (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    release_year int NOT NULL,
    rating numeric(3,1),
    genres text[] NOT NULl,
    director text NOT NULL 
);
