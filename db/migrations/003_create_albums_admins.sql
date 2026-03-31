-- Подключаем расширение pgcrypto, чтобы использовать gen_random_uuid().
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Создаём таблицу администраторов сайта.
CREATE TABLE IF NOT EXISTS admins (

  -- Уникальный идентификатор администратора.
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  -- Email для входа в админку. Должен быть уникальным.
  email TEXT NOT NULL UNIQUE,

  -- Хеш пароля, а не сам пароль.
  password_hash TEXT NOT NULL,

  -- Отображаемое имя администратора.
  display_name TEXT NOT NULL DEFAULT '',

  -- Активен ли аккаунт администратора.
  is_active BOOLEAN NOT NULL DEFAULT TRUE,

  -- Дата создания записи.
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Создаём таблицу альбомов.
CREATE TABLE IF NOT EXISTS albums (

  -- Уникальный идентификатор альбома.
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  -- Человекочитаемый slug для URL, например: easytone-2026
  slug TEXT NOT NULL UNIQUE,

  -- Название альбома.
  title TEXT NOT NULL,

  -- Ссылка на обложку альбома.
  cover_url TEXT NOT NULL DEFAULT '',

  -- Описание альбома.
  description TEXT NOT NULL DEFAULT '',

  -- Публичный ли альбом.
  is_public BOOLEAN NOT NULL DEFAULT TRUE,

  -- Порядок вывода альбомов на сайте.
  sort_order INT NOT NULL DEFAULT 0,

  -- Дата создания альбома.
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Создаём новую версию таблицы фотографий.
-- Назвали её photos_v2, чтобы не ломать старую таблицу photos прямо сейчас.
CREATE TABLE IF NOT EXISTS photos_v2 (

  -- Уникальный идентификатор фотографии.
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  -- ID альбома, к которому относится фото.
  album_id UUID NOT NULL REFERENCES albums(id) ON DELETE CASCADE,

  -- Название фотографии.
  title TEXT NOT NULL DEFAULT '',

  -- Описание фотографии.
  description TEXT NOT NULL DEFAULT '',

  -- Ссылка на полное изображение.
  image_url TEXT NOT NULL,

  -- Ссылка на превью/thumbnail.
  thumb_url TEXT NOT NULL DEFAULT '',

  -- Порядок отображения фото внутри альбома.
  sort_order INT NOT NULL DEFAULT 0,

  -- Дата добавления фотографии.
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Создаём индекс для быстрого поиска фото по album_id.
CREATE INDEX IF NOT EXISTS idx_photos_v2_album_id ON photos_v2(album_id);

-- Создаём индекс для сортировки альбомов на главной.
CREATE INDEX IF NOT EXISTS idx_albums_sort_order ON albums(sort_order);

-- Создаём индекс для сортировки фотографий внутри альбома.
CREATE INDEX IF NOT EXISTS idx_photos_v2_sort_order ON photos_v2(sort_order);