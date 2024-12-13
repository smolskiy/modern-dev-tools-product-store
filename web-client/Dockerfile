# Используем официальный образ Nginx
FROM nginx:alpine

# Удаляем стандартную страницу Nginx
RUN rm -rf /usr/share/nginx/html/*

# Копируем ваш HTML-файл в папку Nginx
COPY client.html /usr/share/nginx/html/index.html

# Открываем порт 80 для веб-сервера
EXPOSE 80

# Запускаем Nginx
CMD ["nginx", "-g", "daemon off;"]
