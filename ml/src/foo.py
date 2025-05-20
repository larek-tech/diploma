import json
import redis

def json_to_redis(source_id, json_file_path, redis_host='localhost', redis_port=6379, redis_db=0):
    r = redis.Redis(host=redis_host, port=redis_port, db=redis_db)
    
    try:
        # Читаем JSON файл
        with open(json_file_path, 'r', encoding='utf-8') as file:
            data = json.load(file)
            
        # Сохраняем каждый элемент в Redis
        for index, item in enumerate(data):
            key = f"{source_id}:{index}"
            r.set(key, json.dumps(item, ensure_ascii=False))
            
        print(f"Успешно сохранено {len(data)} записей в Redis с префиксом '{source_id}'")
    except Exception as e:
        print(f"Ошибка при обработке файла: {e}")
    finally:
        r.close()


if __name__ == "__main__":
    json_to_redis(source_id="d028f055-c743-4b5f-9995-bd9fcf3b4330", json_file_path="/project/src/data.json",redis_host="192.168.1.5")