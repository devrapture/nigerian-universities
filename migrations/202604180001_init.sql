CREATE TABLE "users" (
  "id" uuid NOT NULL,
  "email" text NOT NULL,
  "name" text NOT NULL,
  "avatar_url" text NOT NULL,
  "provider" text NOT NULL,
  "provider_id" text NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX "idx_users_email" ON "users" ("email");
CREATE INDEX "idx_users_deleted_at" ON "users" ("deleted_at");

CREATE TABLE "institutions" (
  "id" uuid NOT NULL,
  "name" text NOT NULL,
  "type" text NOT NULL,
  "vice_chancellor" text NOT NULL,
  "website" text NOT NULL,
  "year_of_establishment" text NOT NULL,
  "last_scraped_at" timestamptz NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);

CREATE INDEX "idx_universities_name" ON "institutions" ("name");
CREATE INDEX "idx_universities_type" ON "institutions" ("type");
CREATE INDEX "idx_universities_vice_chancellor" ON "institutions" ("vice_chancellor");
CREATE INDEX "idx_universities_website" ON "institutions" ("website");
CREATE INDEX "idx_universities_year_of_establishment" ON "institutions" ("year_of_establishment");
CREATE INDEX "idx_institutions_deleted_at" ON "institutions" ("deleted_at");

CREATE TABLE "product_keys" (
  "id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "key_hash" text NOT NULL,
  "key_prefix" text NOT NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "last_used_at" timestamptz NULL,
  "revoked_at" timestamptz NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_product_keys_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);

CREATE UNIQUE INDEX "idx_product_keys_key_hash" ON "product_keys" ("key_hash");
CREATE INDEX "idx_product_keys_user_id" ON "product_keys" ("user_id");
CREATE UNIQUE INDEX "idx_one_active_key_per_user"
  ON "product_keys" ("user_id", "is_active")
  WHERE is_active = true;
CREATE INDEX "idx_product_keys_deleted_at" ON "product_keys" ("deleted_at");
