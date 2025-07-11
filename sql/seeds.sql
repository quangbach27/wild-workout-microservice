-- Insert dates (each represents a full day)
INSERT INTO dates (date) VALUES
  ('2025-06-14 00:00:00'),
  ('2025-06-15 00:00:00');

-- Insert hours (each hour belongs to a date)
INSERT INTO hours (hour, availability, date) VALUES
  -- For 2025-06-14
  ('2025-06-14 12:00:00', 'available', '2025-06-14'),
  ('2025-06-14 13:00:00', 'not_available', '2025-06-14'),
  ('2025-06-14 14:00:00', 'training_scheduled', '2025-06-14'),

  -- For 2025-06-15
  ('2025-06-15 12:00:00', 'available', '2025-06-15'),
  ('2025-06-15 13:00:00', 'available', '2025-06-15'),
  ('2025-06-15 15:00:00', 'not_available', '2025-06-15');
