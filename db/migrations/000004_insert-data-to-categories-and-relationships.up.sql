BEGIN;

INSERT INTO categories (category_id, name, description)
VALUES (1, 'Phone', 'All kinds of phones.'),
       (2, 'Smartphone', 'Touchscreen-based devices running on platforms like Android or iOS.'),
       (3, 'Flip', 'These are feature phones or smartphones that fold in half, with a hinge in the middle.'),
       (4, 'Foldable', 'Advanced smartphones that have a flexible display, allowing the screen to fold.'),
       (5, 'Feature', 'Also known as button phones, they have physical buttons and limited features.'),
       (6, 'Gaming',
        'Specifically optimized for gaming performance with high-refresh-rate screens and powerful hardware.'),
       ---
       (7, 'SoundEquipment', 'All kinds of sound equipments.'),
       (8, 'TWS',
        'Wireless earbuds that come with a charging case. They offer the convenience of being entirely free from any wires.'),
       (9, 'Neckband',
        'Wireless earphones that have a band going around the neck for added security and easier access to controls.'),
       (10, 'Speaker',
        'Portable speakers come in various forms like Bluetooth speakers, smart speakers, and more traditional wired speakers.'),
       (11, 'Wired',
        'These require a physical connection to the audio source via a 3.5mm jack, Type-C or other connectors.'),
       (12, 'Overhead',
        'Large headphones that completely cover the ear, offering better sound isolation and usually better sound quality.'),
       (13, 'IEM',
        'High-quality earphones used by musicians and audiophiles for superior sound isolation and audio fidelity.'),
       ---
       ---
       (14, 'Wearable', 'All kinds of sound wearables.'),
       (15, 'SmartWatch', 'Fitness trackers in band form factor'),
       (16, 'SmartBand', 'Fitness trackers in watch form factor'),
       (17, 'VintageWatch', 'Vintage and classic regular watches')
;

-- restart category id sequence
alter sequence categories_category_id_seq restart 18;

-- Inserting into category_relationships table
INSERT INTO category_relationships (ancestor_id, descendant_id, level)
VALUES
    -- For parent category 'Phone'
    (1, 1, 0),
    (1, 2, 1),
    (1, 3, 1),
    (1, 4, 1),
    (1, 5, 1),
    (1, 6, 1),

    -- For parent category 'SoundEquipment'
    (7, 7, 0),
    (7, 8, 1),
    (7, 9, 1),
    (7, 10, 1),
    (7, 11, 1),
    (7, 12, 1),
    (7, 13, 1),

    -- For parent category 'Wearable'
    (14, 14, 0),
    (14, 15, 1),
    (14, 16, 1),
    (14, 17, 1)
;

COMMIT;
