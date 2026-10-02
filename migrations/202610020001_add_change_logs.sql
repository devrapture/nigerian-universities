CREATE TABLE "change_logs" (
  "id" uuid NOT NULL,
  "institution_id" uuid NOT NULL,
  "institution_name" text NOT NULL,
  "institution_type" text NOT NULL,
  "change_type" text NOT NULL,
  "changes" jsonb NOT NULL DEFAULT '[]'::jsonb,
  "changed_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "change_logs_change_type_check" CHECK ("change_type" IN ('created', 'updated'))
);

CREATE INDEX "idx_change_logs_institution_id" ON "change_logs" ("institution_id");
CREATE INDEX "idx_change_logs_changed_at" ON "change_logs" ("changed_at" DESC);
