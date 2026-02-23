-- USERS
CREATE TABLE users (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) NOT NULL,
    email       VARCHAR(150) UNIQUE NOT NULL,
    password    VARCHAR(255) NOT NULL,        -- bcrypt hashed
    phone       VARCHAR(20),
    role        VARCHAR(20) NOT NULL DEFAULT 'customer', -- customer | restaurant_owner | admin
    is_verified BOOLEAN DEFAULT FALSE,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW()
);

-- RESTAURANTS
CREATE TABLE restaurants (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name        VARCHAR(150) NOT NULL,
    description TEXT,
    address     TEXT NOT NULL,
    city        VARCHAR(100) NOT NULL,
    image_url   TEXT,
    is_active   BOOLEAN DEFAULT TRUE,
    rating      DECIMAL(2,1) DEFAULT 0.0,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW()
);

-- MENU CATEGORIES (e.g. Starters, Main Course, Drinks)
CREATE TABLE menu_categories (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    name          VARCHAR(100) NOT NULL,
    display_order INT DEFAULT 0,
    created_at    TIMESTAMP DEFAULT NOW()
);

-- MENU ITEMS
CREATE TABLE menu_items (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id UUID NOT NULL REFERENCES menu_categories(id) ON DELETE CASCADE,
    name        VARCHAR(150) NOT NULL,
    description TEXT,
    price       DECIMAL(10,2) NOT NULL,
    image_url   TEXT,
    is_available BOOLEAN DEFAULT TRUE,
    is_veg      BOOLEAN DEFAULT FALSE,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW()
);

-- INDEXES for fast lookups
CREATE INDEX idx_restaurants_city ON restaurants(city);
CREATE INDEX idx_restaurants_owner ON restaurants(owner_id);
CREATE INDEX idx_menu_items_category ON menu_items(category_id);
CREATE INDEX idx_users_email ON users(email);