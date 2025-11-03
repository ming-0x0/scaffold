CREATE ROLE scaffold WITH LOGIN PASSWORD 'password';
CREATE DATABASE scaffold
    WITH OWNER scaffold
    ENCODING 'UTF8'
    LC_COLLATE='vi_VN.utf8'
    LC_CTYPE='vi_VN.utf8'
    TEMPLATE template0;

GRANT ALL PRIVILEGES ON DATABASE scaffold TO scaffold;
