CREATE OR REPLACE FUNCTION notify_vocab_update()
RETURNS trigger AS $$
BEGIN
    PERFORM pg_notify(
        'vocab_updates',
        row_to_json(NEW)::text
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER vocab_change_trigger
AFTER INSERT OR UPDATE OR DELETE ON vocabularies
FOR EACH ROW EXECUTE FUNCTION notify_vocab_update();
