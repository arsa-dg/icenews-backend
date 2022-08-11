CREATE OR REPLACE FUNCTION update_count_news() RETURNS TRIGGER AS $$ 
  BEGIN
    UPDATE news
    SET comment = comment+1
    WHERE news.id = NEW.news_id;
	RETURN NULL;
  END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER count_comment
  AFTER INSERT
  ON comments
  FOR EACH ROW
  EXECUTE FUNCTION update_count_news();