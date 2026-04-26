-- Migration: Remove /app prefix from menu routes
UPDATE menus SET route = REPLACE(route, '/app/', '/') WHERE route LIKE '/app/%';
UPDATE menus SET route = '/' WHERE route = '/app';
