---Функция для проставления кода подтверждения (по умолчанию)
CREATE FUNCTION default_sms_code_verification() RETURNS trigger AS $$
BEGIN
    update sessions set verification_code = '7777' where session_id = new.session_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

---Триггер срабатывает при добавление новой сессии и добавляет код подтверждения по умолчанию
CREATE TRIGGER add_default_sms_code_verification  after INSERT ON sessions
    FOR EACH ROW EXECUTE PROCEDURE default_sms_code_verification();