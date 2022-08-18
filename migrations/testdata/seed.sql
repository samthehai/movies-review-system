BEGIN;

SET time_zone = '+00:00';

SELECT 'insert users';

INSERT INTO `users` (`username`, `email`, `hashed_password`)
VALUES (
    'testuser',
    'testuser@gmail.com',
    '$2a$10$qGzkPHjjh/n8N60ARb.BvObjkthrEFF.NCjPKN3RPqDQbpec0JEtG' -- password: secret --
  );

SELECT 'insert movies';

INSERT INTO `movies` (
    `original_title`,
    `original_language`,
    `overview`,
    `adult`,
    `release_date`,
    `budget`,
    `revenue`
  )
VALUES (
    "accumsan sed, facilisis vitae,",
    "Nigeria",
    "risus. Donec nibh enim, gravida sit amet, dapibus id, blandit at, nisi. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Proin vel nisl. Quisque fringilla euismod enim. Etiam gravida molestie arcu. Sed eu nibh vulputate mauris sagittis placerat. Cras dictum ultricies ligula. Nullam enim. Sed nulla ante, iaculis nec, eleifend non, dapibus rutrum, justo. Praesent luctus. Curabitur egestas nunc sed libero. Proin sed turpis nec mauris blandit mattis. Cras",
    true,
    "2022-01-04 02:20:10",
    6036080,
    2576380
  ),
  (
    "vitae velit egestas lacinia.",
    "Singapore",
    "dui. Suspendisse ac metus vitae velit egestas lacinia. Sed congue, elit sed consequat auctor, nunc nulla vulputate dui, nec tempus mauris erat eget ipsum. Suspendisse sagittis. Nullam vitae diam. Proin dolor. Nulla semper tellus id nunc interdum feugiat. Sed nec metus facilisis lorem tristique aliquet. Phasellus fermentum convallis ligula. Donec luctus aliquet odio. Etiam ligula tortor, dictum eu, placerat eget, venenatis a, magna. Lorem",
    true,
    "2022-06-07 16:04:01",
    8028273,
    9315618
  ),
  (
    "nunc sit",
    "Belgium",
    "Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos hymenaeos. Mauris ut quam vel sapien imperdiet ornare. In faucibus. Morbi vehicula. Pellentesque tincidunt tempus risus. Donec egestas. Duis ac arcu. Nunc mauris. Morbi non sapien molestie orci tincidunt adipiscing. Mauris molestie pharetra nibh. Aliquam ornare, libero at auctor ullamcorper, nisl arcu iaculis enim, sit amet ornare lectus justo eu arcu. Morbi sit amet massa. Quisque porttitor eros nec",
    false,
    "2021-12-04 00:59:26",
    432042,
    2997960
  ),
  (
    "semper pretium neque. Morbi quis urna.",
    "Vietnam",
    "Sed pharetra, felis eget varius ultrices, mauris ipsum porta elit, a feugiat tellus lorem eu metus. In lorem. Donec elementum, lorem ut aliquam iaculis, lacus pede sagittis augue, eu tempor erat neque non quam. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Aliquam fringilla cursus purus.",
    true,
    "2023-04-24 14:17:01",
    1086918,
    8711503
  ),
  (
    "arcu. Vivamus sit amet risus. Donec egestas. Aliquam",
    "Belgium",
    "egestas, urna justo faucibus lectus, a sollicitudin orci sem eget massa. Suspendisse eleifend. Cras sed leo. Cras vehicula aliquet libero. Integer in magna. Phasellus dolor elit, pellentesque a, facilisis non, bibendum sed, est. Nunc laoreet lectus quis massa. Mauris vestibulum, neque sed dictum eleifend, nunc risus varius orci, in consequat enim diam vel arcu. Curabitur ut odio vel est tempor bibendum. Donec felis orci, adipiscing non, luctus sit amet, faucibus ut, nulla.",
    false,
    "2023-01-09 02:25:17",
    3568324,
    5716009
  );

COMMIT;
