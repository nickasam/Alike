-- Insert test data for development
INSERT INTO users (phone, nickname, bio, location_lat, location_lng, location_name) VALUES
('+8613800138000', 'Demo User 1', 'Love sports and travel', 31.2304, 121.4737, 'Shanghai'),
('+8613800138001', 'Demo User 2', 'Looking for friends', 31.2204, 121.4637, 'Shanghai'),
('+8613800138002', 'Demo User 3', 'Fitness enthusiast', 31.2404, 121.4837, 'Shanghai');

-- Add tags
INSERT INTO user_tags (user_id, tag) 
SELECT id, 'sports' FROM users WHERE nickname = 'Demo User 1';
INSERT INTO user_tags (user_id, tag) 
SELECT id, 'travel' FROM users WHERE nickname = 'Demo User 1';
INSERT INTO user_tags (user_id, tag) 
SELECT id, 'music' FROM users WHERE nickname = 'Demo User 2';
INSERT INTO user_tags (user_id, tag) 
SELECT id, 'fitness' FROM users WHERE nickname = 'Demo User 3';

-- Add images
INSERT INTO user_images (user_id, image_url, order_index)
SELECT id, 'https://via.placeholder.com/300', 0 FROM users WHERE nickname = 'Demo User 1';
INSERT INTO user_images (user_id, image_url, order_index)
SELECT id, 'https://via.placeholder.com/300', 0 FROM users WHERE nickname = 'Demo User 2';
INSERT INTO user_images (user_id, image_url, order_index)
SELECT id, 'https://via.placeholder.com/300', 0 FROM users WHERE nickname = 'Demo User 3';
