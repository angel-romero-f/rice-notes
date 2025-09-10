DROP TRIGGER IF EXISTS update_notes_updated_at ON notes;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP INDEX IF EXISTS idx_notes_uploaded_at;
DROP INDEX IF EXISTS idx_notes_course_id;
DROP INDEX IF EXISTS idx_notes_user_email;
DROP TABLE IF EXISTS notes;