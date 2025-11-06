-- GECONSOLIDEERDE V11 - CMS/SUPABASE DATA MIGRATIE
-- Dit bestand combineert de logica van de volledige V1_31 reeks (A t/m I).
-- Het maakt alle 8 tabellen aan en voegt de data in in de juiste
-- volgorde om aan FOREIGN KEY constraints te voldoen.

-- Deel A: Photos (van V1_31A)
CREATE TABLE IF NOT EXISTS photos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    url TEXT NOT NULL,
    alt_text TEXT,
    visible BOOLEAN NOT NULL DEFAULT TRUE,
    thumbnail_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    title TEXT,
    description TEXT,
    year INTEGER,
    cloudinary_folder TEXT
);

INSERT INTO photos (id, url, alt_text, visible, thumbnail_url, created_at, updated_at, title, description, year, cloudinary_folder) VALUES
('ee20de98-2fa8-4e23-8bf1-0b705b55aa7c', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1747636922/vv5v84gadrf02rl3iiji.jpg', '1000016660', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747636922/vv5v84gadrf02rl3iiji.jpg', '2025-05-19 06:42:03.254692+00', '2025-05-19 06:42:03.254692+00', '1000016660', null, null, null),
('754ceb60-f4d8-4434-9338-065339e636e8', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513150/d23b6xefqsaxekgnkpeq.jpg', 'WhatsApp Image 2025-05-17 at 20', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513150/d23b6xefqsaxekgnkpeq.jpg', '2025-05-17 20:19:11.491775+00', '2025-05-17 20:19:11.491775+00', 'WhatsApp Image 2025-05-17 at 20', null, null, null),
('eafe6902-3d8e-4929-8cf3-c9b49429adac', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513150/s82ykrgnv8zuwxrraixd.jpg', 'WhatsApp Image 2025-05-17 at 20', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513150/s82ykrgnv8zuwxrraixd.jpg', '2025-05-17 20:19:10.549969+00', '2025-05-17 20:19:10.549969+00', 'WhatsApp Image 2025-05-17 at 20', null, null, null),
('db8dbb96-c528-4ae2-adf7-fc0cd165ad3c', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513149/t2sanzejw8lztqbhbu0h.jpg', 'WhatsApp Image 2025-05-17 at 20', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513149/t2sanzejw8lztqbhbu0h.jpg', '2025-05-17 20:19:09.809387+00', '2025-05-17 20:19:09.809387+00', 'WhatsApp Image 2025-05-17 at 20', null, null, null),
('19312aba-29ab-4169-83eb-934a80763ad8', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513148/c6bzsdf9osgub9cundss.jpg', 'WhatsApp Image 2025-05-17 at 20', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513148/c6bzsdf9osgub9cundss.jpg', '2025-05-17 20:19:08.91355+00', '2025-05-17 20:19:08.91355+00', 'WhatsApp Image 2025-05-17 at 20', null, null, null),
('7eb75484-b2cb-41a7-bdc6-b16a8811ecca', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513147/jwkot927pq1nb8kiwg4h.jpg', 'WhatsApp Image 2025-05-17 at 20', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513147/jwkot927pq1nb8kiwg4h.jpg', '2025-05-17 20:19:08.151916+00', '2025-05-17 20:19:08.151916+00', 'WhatsApp Image 2025-05-17 at 20', null, null, null),
('2debaea6-d5ab-4dd3-ae1b-15919c80245c', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513146/k9o1so9g7jh98dpigakj.jpg', 'WhatsApp Image 2025-05-17 at 20', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513146/k9o1so9g7jh98dpigakj.jpg', '2025-05-17 20:19:07.245565+00', '2025-05-17 20:19:07.245565+00', 'WhatsApp Image 2025-05-17 at 20', null, null, null),
('096c8248-2b05-4337-b256-f857419fa9bd', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513145/cg1my6knmyme8b9ocpre.jpg', 'WhatsApp Image 2025-05-17 at 20', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513145/cg1my6knmyme8b9ocpre.jpg', '2025-05-17 20:19:06.156387+00', '2025-05-17 20:19:06.156387+00', 'WhatsApp Image 2025-05-17 at 20', null, null, null),
('7d525fa5-192a-4b5f-9f03-7b8f9d4eb038', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513144/tlhlt2etlp0wd3hhh1e1.jpg', 'WhatsApp Image 2025-05-17 at 20', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513144/tlhlt2etlp0wd3hhh1e1.jpg', '2025-05-17 20:19:05.098892+00', '2025-05-17 20:19:05.098892+00', 'WhatsApp Image 2025-05-17 at 20', null, null, null),
('9db233dc-58ff-4352-be6e-fb0d3b30798a', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513143/wmtpajgjjvlf7gdpromp.jpg', 'WhatsApp Image 2025-05-17 at 20', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513143/wmtpajgjjvlf7gdpromp.jpg', '2025-05-17 20:19:04.357247+00', '2025-05-17 20:19:04.357247+00', 'WhatsApp Image 2025-05-17 at 20', null, null, null),
('fded98ba-c7fa-4c42-ad2e-6c7afea5b8fc', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513143/jtcxza8j43cwwyynrq5g.jpg', 'WhatsApp Image 2025-05-17 at 19', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513143/jtcxza8j43cwwyynrq5g.jpg', '2025-05-17 20:19:03.527768+00', '2025-05-17 20:19:03.527768+00', 'WhatsApp Image 2025-05-17 at 19', null, null, null),
('c5cc7e20-9e4d-4e09-bcf5-0c2edcb83360', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513142/w9nvohxsntxfiy4cevbf.jpg', 'WhatsApp Image 2025-05-17 at 19', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513142/w9nvohxsntxfiy4cevbf.jpg', '2025-05-17 20:19:02.830463+00', '2025-05-17 20:19:02.830463+00', 'WhatsApp Image 2025-05-17 at 19', null, null, null),
('ac4fdceb-3598-4ba9-9d61-4f40ef0e44e3', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513141/dol0bmbwzaamhvf7zeni.jpg', 'WhatsApp Image 2025-05-17 at 19', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513141/dol0bmbwzaamhvf7zeni.jpg', '2025-05-17 20:19:02.129184+00', '2025-05-17 20:19:02.129184+00', 'WhatsApp Image 2025-05-17 at 19', null, null, null),
('be16d3ae-c8a4-457b-9806-fee919bc79a2', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513140/grdni6fzojt466urmkya.jpg', 'WhatsApp Image 2025-05-17 at 19', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513140/grdni6fzojt466urmkya.jpg', '2025-05-17 20:19:01.404594+00', '2025-05-17 20:19:01.404594+00', 'WhatsApp Image 2025-05-17 at 19', null, null, null),
('c199f984-64e5-4405-b23c-e6ff4a3eaed3', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513140/men0m6mk5f505uoaclhf.jpg', 'WhatsApp Image 2025-05-17 at 19', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513140/men0m6mk5f505uoaclhf.jpg', '2025-05-17 20:19:00.628702+00', '2025-05-17 20:19:00.628702+00', 'WhatsApp Image 2025-05-17 at 19', null, null, null),
('e971f01a-1d4a-4400-bb54-3428fbc69a98', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513139/j19kt4rorb9ybtpx0x1r.jpg', 'WhatsApp Image 2025-05-17 at 19', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513139/j19kt4rorb9ybtpx0x1r.jpg', '2025-05-17 20:18:59.94824+00', '2025-05-17 20:18:59.94824+00', 'WhatsApp Image 2025-05-17 at 19', null, null, null),
('991bfd41-1365-4bb4-9ac8-8519fc9bfb32', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513138/uelmhlfiuccmv2slqbaw.jpg', 'WhatsApp Image 2025-05-17 at 19', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513138/uelmhlfiuccmv2slqbaw.jpg', '2025-05-17 20:18:59.143222+00', '2025-05-17 20:18:59.143222+00', 'WhatsApp Image 2025-05-17 at 19', null, null, null),
('18eb6cd7-b9a4-4b21-8363-ef77f420ac09', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513137/zgflyhmixak9ci0warv4.jpg', 'WhatsApp Image 2025-05-17 at 19', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513137/zgflyhmixak9ci0warv4.jpg', '2025-05-17 20:18:58.33864+00', '2025-05-17 20:18:58.33864+00', 'WhatsApp Image 2025-05-17 at 19', null, null, null),
('83e75346-54b4-40b7-8ba4-2d43b3d8c867', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513136/iimq27dhkyimotercqan.jpg', 'WhatsApp Image 2025-05-17 at 19', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513136/iimq27dhkyimotercqan.jpg', '2025-05-17 20:18:57.33427+00', '2025-05-17 20:18:57.33427+00', 'WhatsApp Image 2025-05-17 at 19', null, null, null),
('49f983ab-ec63-4489-93ce-ba9272ba7f49', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513135/bwpkhyzltxncxkza2lsz.jpg', 'WhatsApp Image 2025-05-17 at 19', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513135/bwpkhyzltxncxkza2lsz.jpg', '2025-05-17 20:18:56.358544+00', '2025-05-17 20:18:56.358544+00', 'WhatsApp Image 2025-05-17 at 19', null, null, null),
('6649c67b-2eb5-4d29-9e1c-0373ea3b1771', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513134/l6o1ibr7lzzbn7glpjcu.jpg', 'WhatsApp Image 2025-05-17 at 19', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513134/l6o1ibr7lzzbn7glpjcu.jpg', '2025-05-17 20:18:55.523962+00', '2025-05-17 20:18:55.523962+00', 'WhatsApp Image 2025-05-17 at 19', null, null, null),
('ec209a88-3f3b-4168-a4c9-1c9275edcffb', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513133/khtzc08kc7wgkta5rh7s.jpg', 'WhatsApp Image 2025-05-17 at 19', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513133/khtzc08kc7wgkta5rh7s.jpg', '2025-05-17 20:18:54.564013+00', '2025-05-17 20:18:54.564013+00', 'WhatsApp Image 2025-05-17 at 19', null, null, null),
('c7ddfbd8-6bf4-4bd9-9285-f6492b98bef4', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513133/unhly8fepi83vegupc6a.jpg', 'WhatsApp Image 2025-05-17 at 19', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513133/unhly8fepi83vegupc6a.jpg', '2025-05-17 20:18:53.737493+00', '2025-05-17 20:18:53.737493+00', 'WhatsApp Image 2025-05-17 at 19', null, null, null),
('10fff5f8-2701-4f34-86f1-b063b252f35a', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513132/zoqmk50gcxuqkda0sqoe.jpg', 'WhatsApp Image 2025-05-17 at 19', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513132/zoqmk50gcxuqkda0sqoe.jpg', '2025-05-17 20:18:53.033148+00', '2025-05-17 20:18:53.033148+00', 'WhatsApp Image 2025-05-17 at 19', null, null, null),
('76278110-be4c-4ffc-bc35-e2d4b14fa42f', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513131/hfo3n8mzetzeqefr418d.jpg', 'WhatsApp Image 2025-05-17 at 19', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513131/hfo3n8mzetzeqefr418d.jpg', '2025-05-17 20:18:52.092844+00', '2025-05-17 20:18:52.092844+00', 'WhatsApp Image 2025-05-17 at 19', null, null, null),
('baf57b95-0ab7-4217-af98-1453a9e6a938', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513130/tbsbwsnimr5fiv1i7pkn.jpg', 'WhatsApp Image 2025-05-17 at 19', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513130/tbsbwsnimr5fiv1i7pkn.jpg', '2025-05-17 20:18:51.341415+00', '2025-05-17 20:18:51.341415+00', 'WhatsApp Image 2025-05-17 at 19', null, null, null),
('0911da4d-6169-4383-962b-9a640b5eac0e', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513129/p8aoklpqxch3jlbukl3r.jpg', 'WhatsApp Image 2025-05-17 at 19', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513129/p8aoklpqxch3jlbukl3r.jpg', '2025-05-17 20:18:50.351733+00', '2025-05-17 20:18:50.351733+00', 'WhatsApp Image 2025-05-17 at 19', null, null, null),
('1350bb1c-b132-4250-82b8-58efb05fb53b', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513128/thtibyrsflsuen2lotv5.jpg', 'WhatsApp Image 2025-05-17 at 19', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513128/thtibyrsflsuen2lotv5.jpg', '2025-05-17 20:18:49.478846+00', '2025-05-17 20:18:49.478846+00', 'WhatsApp Image 2025-05-17 at 19', null, null, null),
('08362e92-340a-432a-b306-153ad27ee686', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513127/thhu8mxkqhfhxjydi5ai.jpg', 'WhatsApp Image 2025-05-17 at 19', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513127/thhu8mxkqhfhxjydi5ai.jpg', '2025-05-17 20:18:48.433401+00', '2025-05-17 20:18:48.433401+00', 'WhatsApp Image 2025-05-17 at 19', null, null, null),
('7c296b82-d35d-4065-8523-dffc976688f7', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513126/j5l0h7pakhadhwcrxo8o.jpg', 'WhatsApp Image 2025-05-17 at 19', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513126/j5l0h7pakhadhwcrxo8o.jpg', '2025-05-17 20:18:47.544035+00', '2025-05-17 20:18:47.544035+00', 'WhatsApp Image 2025-05-17 at 19', null, null, null),
('5b38433f-3c33-49cb-950a-e0bc7eb21a80', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1745793210/rirkdj5bav7k0pvvtfoq.jpg', '1000014905', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1745793210/rirkdj5bav7k0pvvtfoq.jpg', '2025-04-27 22:33:31.229221+00', '2025-04-27 22:33:31.229221+00', '1000014905', null, null, null),
('dbb68d91-3f6f-46c3-b751-a13f34269b79', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1734808286/ubjiz1fal82jh42nzjmg.jpg', 'ubjiz1fal82jh42nzjmg', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1734808286/ubjiz1fal82jh42nzjmg.jpg', '2025-04-19 00:30:16.928614+00', '2025-04-20 02:08:06.401666+00', '321', null, 2025, null),
('f4ce7f8e-b573-4602-9e12-69e0585df779', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1745022279/dlrhmdl4gcddunkqzkei.jpg', 'dlrhmdl4gcddunkqzkei', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1745022279/dlrhmdl4gcddunkqzkei.jpg', '2025-04-19 00:30:16.928614+00', '2025-04-20 02:07:54.136061+00', '123', null, 2025, null),
('e7b84300-7158-475b-a79b-10e97b416d58', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1745022320/k7qobfrinlsjrqrxxdfy.jpg', 'k7qobfrinlsjrqrxxdfy', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1745022320/k7qobfrinlsjrqrxxdfy.jpg', '2025-04-19 00:30:16.928614+00', '2025-04-20 02:08:27.195419+00', '231', null, 2025, null),
('26188a5b-d542-4674-8ae9-b63520fbd4b2', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1734808313/bmkf9wcfwrseamcgdu9t.jpg', 'bmkf9wcfwrseamcgdu9t', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1734808313/bmkf9wcfwrseamcgdu9t.jpg', '2025-04-19 00:30:16.928614+00', '2025-04-20 02:08:16.199794+00', '213', null, 2025, null),
('d78d9d0f-29ac-42a1-951a-a0bcb355d227', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1739543787/vdohoeldmm6iiv9cwikm.jpg', '1000011474', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1739543787/vdohoeldmm6iiv9cwikm.jpg', '2025-02-14 14:36:28.503034+00', '2025-04-20 02:09:25.344705+00', '011474', null, 2025, null),
('af027789-d718-4899-88b2-d8f0814b73ee', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1733163484/671433bb6de927107d736c2a_DKL19-p-800_slbxvx.webp', 'Koninklijke Loop 2023 - Sfeerimpressie', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200/v1733163484/671433bb6de927107d736c2a_DKL19-p-800_slbxvx.webp', '2024-12-22 20:03:33.896605+00', '2025-04-17 19:05:18.330672+00', 'Koninklijke Loop 2023 - Sfeerimpressie', null, 2024, null),
('6fd90008-f1ca-4e3d-9915-32b914208239', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1733163485/6714330a3210f9f512505de3_dkl2-p-800_pvztxb.webp', 'Koninklijke Loop 2023 - Parcours impressie', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200/v1733163485/6714330a3210f9f512505de3_dkl2-p-800_pvztxb.webp', '2024-12-22 20:03:33.896605+00', '2025-03-17 19:50:37.910515+00', 'Koninklijke Loop 2023 - Parcours impressie', null, 2024, null),
('3acaeca2-aa51-4cd6-9fad-7b20b607cba7', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1733163486/6714336722ccd9c2f9c73177_DKL12-p-800_hguzw1.webp', 'Koninklijke Loop 2023 - Evenement overzicht', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200/v1733163486/6714336722ccd9c2f9c73177_DKL12-p-800_hguzw1.webp', '2024-12-22 20:03:33.896605+00', '2025-03-17 19:50:27.251872+00', 'Koninklijke Loop 2023 - Evenement overzicht', null, 2024, null),
('4059d3cf-0e92-44c2-bbfa-93235dc19eae', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1733163486/67143385825a5682baca0273_DKL13-p-800_m9sc65.webp', 'Koninklijke Loop 2023 - Sfeerbeeld', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200/v1733163486/67143385825a5682baca0273_DKL13-p-800_m9sc65.webp', '2024-12-22 20:03:33.896605+00', '2025-03-17 19:50:30.097262+00', 'Koninklijke Loop 2023 - Sfeerbeeld', null, 2024, null),
('245ddf1a-61ab-4b64-b8b2-13d5079d6592', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1733163485/6714335d14c39bfea035fa1b_DKL9-p-800_umfakb.webp', 'Koninklijke Loop 2023 - Finish moment', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200/v1733163485/6714335d14c39bfea035fa1b_DKL9-p-800_umfakb.webp', '2024-12-22 20:03:33.896605+00', '2025-03-17 19:50:24.427491+00', 'Koninklijke Loop 2023 - Finish moment', null, 2024, null),
('0334c63c-230a-46c7-b87d-3f4a5cc946c0', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1733163485/6714333d9d3ac37ffaa1d539_DKL6-p-800_o16v1m.webp', 'Koninklijke Loop 2023 - Deelnemers samen', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200/v1733163485/6714333d9d3ac37ffaa1d539_DKL6-p-800_o16v1m.webp', '2024-12-22 20:03:33.896605+00', '2025-03-17 19:50:19.507503+00', 'Koninklijke Loop 2023 - Deelnemers samen', '', 2024, null),
('e1350a21-3d42-4252-bc56-5d9ce264ff41', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1733163485/67143326f367de40be79ebfa_DKL4-p-800_ceagh7.webp', 'Koninklijke Loop 2023 - Lopers onderweg', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200/v1733163485/67143326f367de40be79ebfa_DKL4-p-800_ceagh7.webp', '2024-12-22 20:03:33.896605+00', '2025-03-17 19:50:56.31221+00', 'Koninklijke Loop 2023 - Lopers onderweg', null, 2024, null),
('8a4d5c20-ea73-4336-8c6b-af7a197ef7c2', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1733163484/671433481981522775f532d1_DKL7-p-800_zdihdc.webp', 'Koninklijke Loop 2023 - Deelnemers in actie', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200/v1733163484/671433481981522775f532d1_DKL7-p-800_zdihdc.webp', '2024-12-22 20:03:33.896605+00', '2025-03-17 19:50:44.13556+00', 'Koninklijke Loop 2023 - Deelnemers in actie', null, 2024, null),
('7431c540-2aa5-42db-b5c5-df55a17856ae', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1733163485/67143316c21bdbd9d9800bcf_DKL3-p-800_kcx0b3.webp', 'Koninklijke Loop 2023 - Groepsfoto', true, 'https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200/v1733163485/67143316c21bdbd9d9800bcf_DKL3-p-800_kcx0b3.webp', '2024-12-22 20:03:33.896605+00', '2025-03-17 19:50:41.386987+00', 'Koninklijke Loop 2023 - Groepsfoto', null, 2024, null)
ON CONFLICT (id) DO NOTHING;

-- Deel B: Albums (van V1_31B)
-- Afhankelijk van 'photos' (voor cover_photo_id)
CREATE TABLE IF NOT EXISTS albums (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    description TEXT,
    cover_photo_id UUID REFERENCES photos(id),
    visible BOOLEAN NOT NULL DEFAULT TRUE,
    order_number INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

INSERT INTO albums (id, title, description, cover_photo_id, visible, order_number, created_at, updated_at) VALUES
('72831c18-4c6c-4dc8-9a2a-ace696e8996b', 'Voorbereidingen ', 'Voor werk', 'dbb68d91-3f6f-46c3-b751-a13f34269b79', true, 3, '2025-02-14 14:35:32.369+00', '2025-02-14 14:35:32.37+00'),
('ce8df963-f118-4296-9f5c-e33308dc7bfa', 'DKL-2024', 'DKL 2024', '8a4d5c20-ea73-4336-8c6b-af7a197ef7c2', true, 2, '2024-12-26 15:17:54.268+00', '2025-02-17 13:32:50.793+00'),
('d51cff45-b958-4370-a983-51e650ffa43e', 'DKL 2025', 'De koninklijke Loop 2025!', '08362e92-340a-432a-b306-153ad27ee686', true, 1, '2025-05-17 20:10:00+00', '2025-05-17 20:10:06.643082+00')
ON CONFLICT (id) DO NOTHING;

-- Deel C: Album Photos (van V1_31C)
-- Afhankelijk van 'photos' en 'albums'
CREATE TABLE IF NOT EXISTS album_photos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    album_id UUID NOT NULL REFERENCES albums(id) ON DELETE CASCADE,
    photo_id UUID NOT NULL REFERENCES photos(id) ON DELETE CASCADE,
    order_number INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

INSERT INTO album_photos (id, album_id, photo_id, order_number, created_at) VALUES
('10592b5a-9f16-48d2-8e95-150e5c5c53e5', 'd51cff45-b958-4370-a983-51e650ffa43e', 'ee20de98-2fa8-4e23-8bf1-0b705b55aa7c', 30, '2025-05-19 06:42:03.513117+00'),
('aa902fe8-9ef7-44db-bbe1-1a48c203a26c', 'd51cff45-b958-4370-a983-51e650ffa43e', 'ec209a88-3f3b-4168-a4c9-1c9275edcffb', 9, '2025-05-17 20:19:11.664743+00'),
('e3734c7c-a99d-421f-9a0c-ac02f08cbaf8', 'd51cff45-b958-4370-a983-51e650ffa43e', 'be16d3ae-c8a4-457b-9806-fee919bc79a2', 17, '2025-05-17 20:19:11.664743+00'),
('1fa5d9f4-942d-4465-ad93-77ed1098d1dd', 'd51cff45-b958-4370-a983-51e650ffa43e', 'ac4fdceb-3598-4ba9-9d61-4f40ef0e44e3', 18, '2025-05-17 20:19:11.664743+00'),
('c80655d5-0091-40bc-95cb-b1fcdab0112e', 'd51cff45-b958-4370-a983-51e650ffa43e', 'c5cc7e20-9e4d-4e09-bcf5-0c2edcb83360', 19, '2025-05-17 20:19:11.664743+00'),
('c64d8242-6200-4c19-b55e-2eb7ca8d13d9', 'd51cff45-b958-4370-a983-51e650ffa43e', 'fded98ba-c7fa-4c42-ad2e-6c7afea5b8fc', 20, '2025-05-17 20:19:11.664743+00'),
('486dcb08-ba33-47a7-af91-612d56e98639', 'd51cff45-b958-4370-a983-51e650ffa43e', '9db233dc-58ff-4352-be6e-fb0d3b30798a', 21, '2025-05-17 20:19:11.664743+00'),
('e7de00f4-fc44-4ff8-abd7-4d9c079c97fb', 'd51cff45-b958-4370-a983-51e650ffa43e', '7d525fa5-192a-4b5f-9f03-7b8f9d4eb038', 22, '2025-05-17 20:19:11.664743+00'),
('74069e01-a97a-45a0-a82e-4c3baf231e23', 'd51cff45-b958-4370-a983-51e650ffa43e', '096c8248-2b05-4337-b256-f857419fa9bd', 23, '2025-05-17 20:19:11.664743+00'),
('ea79a259-5dee-4f64-8009-00ddeae4c636', 'd51cff45-b958-4370-a983-51e650ffa43e', '2debaea6-d5ab-4dd3-ae1b-15919c80245c', 24, '2025-05-17 20:19:11.664743+00'),
('09101c7c-e7ac-4f99-af78-63d06f343d94', 'd51cff45-b958-4370-a983-51e650ffa43e', '7eb75484-b2cb-41a7-bdc6-b16a8811ecca', 25, '2025-05-17 20:19:11.664743+00'),
('65f4c380-e20c-42da-8a06-70d85d3f3ec6', 'd51cff45-b958-4370-a983-51e650ffa43e', '19312aba-29ab-4169-83eb-934a80763ad8', 26, '2025-05-17 20:19:11.664743+00'),
('bf01850a-e7d1-43c2-a614-7e1a65f3a8fc', 'd51cff45-b958-4370-a983-51e650ffa43e', 'db8dbb96-c528-4ae2-adf7-fc0cd165ad3c', 27, '2025-05-17 20:19:11.664743+00'),
('df7ce81b-0480-4955-ad42-cf17dcef84ab', 'd51cff45-b958-4370-a983-51e650ffa43e', 'eafe6902-3d8e-4929-8cf3-c9b49429adac', 28, '2025-05-17 20:19:11.664743+00'),
('898ea918-4be5-4c7f-8e6b-5d8b73198974', 'd51cff45-b958-4370-a983-51e650ffa43e', '754ceb60-f4d8-4434-9338-065339e636e8', 29, '2025-05-17 20:19:11.664743+00'),
('ae9666e2-5867-4d92-b25b-b6e665414dc7', 'd51cff45-b958-4370-a983-51e650ffa43e', '7c296b82-d35d-4065-8523-dffc976688f7', 2, '2025-05-17 20:19:11.664743+00'),
('9ce0e5fc-b7ea-4f63-bf63-c5c0f720ed61', 'd51cff45-b958-4370-a983-51e650ffa43e', '08362e92-340a-432a-b306-153ad27ee686', 1, '2025-05-17 20:19:11.664743+00'),
('e519317a-c8a1-47ff-a935-fc4e4c8c6afc', 'd51cff45-b958-4370-a983-51e650ffa43e', '1350bb1c-b132-4250-82b8-58efb05fb53b', 3, '2025-05-17 20:19:11.664743+00'),
('28974549-324c-49ff-b441-82fe90c947cd', 'd51cff45-b958-4370-a983-51e650ffa43e', '0911da4d-6169-4383-962b-9a640b5eac0e', 4, '2025-05-17 20:19:11.664743+00'),
('d2f53994-5156-4fde-a63c-02e6c39c198e', 'd51cff45-b958-4370-a983-51e650ffa43e', 'baf57b95-0ab7-4217-af98-1453a9e6a938', 5, '2025-05-17 20:19:11.664743+00'),
('0bf73fa0-f48d-4272-87aa-abb53a5afd5a', 'd51cff45-b958-4370-a983-51e650ffa43e', '76278110-be4c-4ffc-bc35-e2d4b14fa42f', 6, '2025-05-17 20:19:11.664743+00'),
('c7cad470-1808-46ab-99af-9b5681c8ca69', 'd51cff45-b958-4370-a983-51e650ffa43e', '10fff5f8-2701-4f34-86f1-b063b252f35a', 7, '2025-05-17 20:19:11.664743+00'),
('f5e86512-f99d-4700-bc32-c2f5b235a670', 'd51cff45-b958-4370-a983-51e650ffa43e', 'c7ddfbd8-6bf4-4bd9-9285-f6492b98bef4', 8, '2025-05-17 20:19:11.664743+00'),
('7551cf97-06fe-4651-ab0b-538dab99f5c6', 'd51cff45-b958-4370-a983-51e650ffa43e', '6649c67b-2eb5-4d29-9e1c-0373ea3b1771', 10, '2025-05-17 20:19:11.664743+00'),
('4e9db229-bcff-4118-b846-3931301fc20c', 'd51cff45-b958-4370-a983-51e650ffa43e', '49f983ab-ec63-4489-93ce-ba9272ba7f49', 11, '2025-05-17 20:19:11.664743+00'),
('78872e14-ad35-4b12-92a7-8e35f59c448b', 'd51cff45-b958-4370-a983-51e650ffa43e', '83e75346-54b4-40b7-8ba4-2d43b3d8c867', 12, '2025-05-17 20:19:11.664743+00'),
('6fbfc586-b7ac-4af6-9053-13ee5cdc4f86', 'd51cff45-b958-4370-a983-51e650ffa43e', '18eb6cd7-b9a4-4b21-8363-ef77f420ac09', 13, '2025-05-17 20:19:11.664743+00'),
('17bf4b95-a2e2-46f9-830d-c86825aaddcc', 'd51cff45-b958-4370-a983-51e650ffa43e', '991bfd41-1365-4bb4-9ac8-8519fc9bfb32', 14, '2025-05-17 20:19:11.664743+00'),
('9e8ef436-84a7-4b42-b5a0-da23dd2caed6', 'd51cff45-b958-4370-a983-51e650ffa43e', 'e971f01a-1d4a-4400-bb54-3428fbc69a98', 15, '2025-05-17 20:19:11.664743+00'),
('622a17b5-98a4-4655-8c7b-db075c1a7d2f', 'd51cff45-b958-4370-a983-51e650ffa43e', 'c199f984-64e5-4405-b23c-e6ff4a3eaed3', 16, '2025-05-17 20:19:11.664743+00'),
('cb67b8ed-cb18-4266-b5dc-d8a8c7b0128b', '72831c18-4c6c-4dc8-9a2a-ace696e8996b', '26188a5b-d542-4674-8ae9-b63520fbd4b2', 15, '2025-04-19 00:31:58.466572+00'),
('5db43327-8974-48c7-accc-33affa2bf999', '72831c18-4c6c-4dc8-9a2a-ace696e8996b', 'dbb68d91-3f6f-46c3-b751-a13f34269b79', 16, '2025-04-19 00:31:58.466572+00'),
('29bea012-0b00-408d-8ff2-beffcd286730', '72831c18-4c6c-4dc8-9a2a-ace696e8996b', 'e7b84300-7158-475b-a79b-10e97b416d58', 14, '2025-04-19 00:31:58.466572+00'),
('47bee588-1be9-43c5-afea-80d67e97a4c2', '72831c18-4c6c-4dc8-9a2a-ace696e8996b', 'f4ce7f8e-b573-4602-9e12-69e0585df779', 13, '2025-04-19 00:31:58.466572+00'),
('fb9d85b1-c9aa-4935-8c1d-848603223d5c', 'ce8df963-f118-4296-9f5c-e33308dc7bfa', 'd78d9d0f-29ac-42a1-951a-a0bcb355d227', 10, '2025-04-18 21:22:50.820523+00'),
('c3035d06-0fff-421a-894b-0b8c804906ae', 'ce8df963-f118-4296-9f5c-e33308dc7bfa', '245ddf1a-61ab-4b64-b8b2-13d5079d6592', 5, '2025-03-17 19:49:06.359018+00'),
('6f2c18d9-b0aa-45d7-bf24-62aa24b4f309', 'ce8df963-f118-4296-9f5c-e33308dc7bfa', '3acaeca2-aa51-4cd6-9fad-7b20b607cba7', 9, '2025-03-17 19:49:06.359018+00'),
('d0a158e2-c89e-4866-9ccb-5493e5cdc445', 'ce8df963-f118-4296-9f5c-e33308dc7bfa', '0334c63c-230a-46c7-b87d-3f4a5cc946c0', 8, '2025-03-17 19:49:06.359018+00'),
('cfecedf5-7b64-4229-80ff-90dc3dc9901a', 'ce8df963-f118-4296-9f5c-e33308dc7bfa', '4059d3cf-0e92-44c2-bbfa-93235dc19eae', 7, '2025-03-17 19:49:06.359018+00'),
('ec507140-12fe-4d96-9b41-78da5b6d5ba0', 'ce8df963-f118-4296-9f5c-e33308dc7bfa', '6fd90008-f1ca-4e3d-9915-32b914208239', 6, '2025-03-17 19:49:06.359018+00'),
('c70e5b3f-197f-4fe1-86bb-cd3a2c5f4b65', 'ce8df963-f118-4296-9f5c-e33308dc7bfa', '8a4d5c20-ea73-4336-8c6b-af7a197ef7c2', 1, '2025-03-17 19:49:06.359018+00'),
('3227ab06-4171-433c-97db-727a0a59b873', 'ce8df963-f118-4296-9f5c-e33308dc7bfa', 'af027789-d718-4899-88b2-d8f0814b73ee', 2, '2025-03-17 19:49:06.359018+00'),
('70055993-7e41-49c1-903e-ff69c8f4fee3', 'ce8df963-f118-4296-9f5c-e33308dc7bfa', 'e1350a21-3d42-4252-bc56-5d9ce264ff41', 3, '2025-03-17 19:49:06.359018+00'),
('8c64e44a-37f8-4fb9-b1af-8195b6d2bf8e', 'ce8df963-f118-4296-9f5c-e33308dc7bfa', '7431c540-2aa5-42db-b5c5-df55a17856ae', 4, '2025-03-17 19:49:06.359018+00'),
('3d958630-2821-4633-8126-23fc9f881863', '72831c18-4c6c-4dc8-9a2a-ace696e8996b', 'd78d9d0f-29ac-42a1-951a-a0bcb355d227', 12, '2025-02-14 14:36:41.258168+00')
ON CONFLICT (id) DO NOTHING;

-- Deel D: Videos (van V1_31D)
CREATE TABLE IF NOT EXISTS videos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    video_id TEXT NOT NULL,
    url TEXT NOT NULL,
    title TEXT,
    description TEXT,
    thumbnail_url TEXT,
    visible BOOLEAN NOT NULL DEFAULT TRUE,
    order_number INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

INSERT INTO videos (id, video_id, url, title, description, thumbnail_url, visible, order_number, created_at, updated_at) VALUES
('14ee164b-50e3-4f59-a3e9-3a54312af9cd', 'q9ngqu', 'https://streamable.com/e/q9ngqu', 'De koninklijkeloop!', 'Preview!!!', null, true, 1, '2025-03-28 21:45:53+00', '2025-04-21 10:49:03.744+00'),
('18d951d2-f5d1-4b6e-95af-a72ac5ff18ff', 'x8zj4k', 'https://streamable.com/e/x8zj4k', 'Promotie De Koninklijke Loop: Flyers verspreiden', 'Bekijk hoe vrijwilligers flyers uitdelen om mensen uit te nodigen voor het DKL wandelevenement.', null, true, 4, '2025-03-03 20:25:54.251108+00', '2025-03-03 20:25:54.251108+00'),
('87502f84-91db-419f-9766-4071ade3e94f', 'tt6k80', 'https://streamable.com/e/0o2qf9', 'Highlights Koninklijke Loop 2024 (Wandelevenement Apeldoorn)', 'Herbeleef de mooiste momenten en de sfeer van De Koninklijke Loop 2024 in deze highlight video.', null, true, 2, '2024-12-22 20:23:06.654419+00', '2024-12-26 23:25:42.699+00'),
('99bbe55b-32ef-46ab-bb59-860ce92f1d58', 'cvfrpi', 'https://streamable.com/e/cvfrpi', 'De spannende start van de Koninklijke Loop 2024', 'Bekijk de start van de deelnemers aan de sponsorloop De Koninklijke Loop 2024.', null, true, 3, '2024-12-22 20:23:06.654419+00', '2024-12-26 23:25:43.979+00'),
('ac987839-edc1-468f-9080-064f894b3e5d', 'tt6k80', 'https://streamable.com/e/tt6k80', 'Koninklijke Loop 2024 - Hoofdevenement', 'Een sfeerimpressie van het hoofdevenement van de Koninklijke Loop 2024, met deelnemers, vrijwilligers en muziek.', null, true, 5, '2024-12-22 20:23:06.654419+00', '2024-12-26 23:25:41.341+00')
ON CONFLICT (id) DO NOTHING;

-- Deel E: Sponsors (van V1_31E)
CREATE TABLE IF NOT EXISTS sponsors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    logo_url TEXT,
    website_url TEXT,
    order_number INTEGER,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    visible BOOLEAN NOT NULL DEFAULT TRUE
);

INSERT INTO sponsors (id, name, description, logo_url, website_url, order_number, is_active, created_at, updated_at, visible) VALUES
('484576a1-2a60-4201-b582-e1f3754ab12a', '3x3 Anders', '3x3 Anders is een zorgbemiddelingsbureau gespecialiseerd in het matchen van zorgaanbieders met gekwalificeerde zorgprofessionals.', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1733166671/3x3anderslogo_itwm3g.webp', 'https://3x3anders.nl/', 4, true, '2024-11-29 09:49:35+00', '2025-04-30 13:07:58.478643+00', true),
('6408a640-cca1-4aaa-b845-04888f62ccec', 'Sterk In Vloeren', 'De website van Sterk In Vloeren biedt een uitgebreid assortiment aan vloeren, waaronder laminaat, PVC-vloeren en tapijt. Ze benadrukken heldere afspraken en hanteren all-in prijzen.', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1733166669/SterkinVloerenLOGO_zrdofb.webp', 'https://sterkinvloeren.nl/', 1, true, '2024-11-29 09:49:35+00', '2025-04-30 13:07:58.478643+00', true),
('6acd1b1e-8fed-4c8b-89cb-85eee9053536', 'Beeldpakker', 'Johan Groot Jebbink, een fotograaf met meer dan tien jaar ervaring, gespecialiseerd in portretfotografie. Actief in Ermelo en internationaal.', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1733166670/BeeldpakkerLogo_wijjmq.webp', 'https://beeldpakker.nl/', 2, true, '2024-11-29 09:49:35+00', '2025-04-30 13:07:58.478643+00', true),
('88cec1c1-3f08-4a6f-9623-c71605fe35b5', 'Bas Visual Story Telling', 'BAS Visual Storytelling, heeft een passie voor content en verhalen. Mijn hobby is uitgegroeid tot een eigen onderneming in het vastleggen van verhalen. Bij BAS Visual Storytelling laten we verhalen niet verstoffen op de plank, maar brengen ze tot leven! Waar ik ga of sta, mijn camera''s gaan met mij mee, leg de mooiste beelden haarscherp vast en breng jouw verhaal tot leven. Dus vertel eens, ''wat is jouw verhaal?''

', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1746017513/krqjbwwerv9hs6hyrhcy.png', 'https://basvisualstorytelling.nl/', 5, true, '2025-04-30 12:51:53.997413+00', '2025-05-01 09:30:16.725422+00', true),
('caa59f1f-65b4-442f-84d2-22cc52212dea', 'Mojo Dojo', 'Mojo Dojo is een veelzijdige studio in Rotterdam die diensten aanbiedt voor creatieve producties, waaronder muziekopnames, podcasts en livestreams.', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1733166669/LogoLayout_1_iphclc.webp', 'https://mojodojo.studio/', 3, true, '2024-11-29 09:49:35+00', '2025-04-30 13:07:58.478643+00', true)
ON CONFLICT (id) DO NOTHING;

-- Deel F: Program Schedule (van V1_31F)
CREATE TABLE IF NOT EXISTS program_schedule (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    time TEXT NOT NULL,
    event_description TEXT NOT NULL,
    category TEXT,
    icon_name TEXT,
    order_number INTEGER,
    visible BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8)
);

INSERT INTO program_schedule (id, time, event_description, category, icon_name, order_number, visible, created_at, updated_at, latitude, longitude) VALUES
('075095c7-925d-411e-bebc-a7fc96a3000a', '12:00u', 'Aanvang Deelnemers 10km bij het coördinatiepunt', 'Aanvang', 'aanvang', 50, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', null, null),
('1302ac11-1b7a-4738-a2dc-535990a09e69', '10:15u', 'Aanvang Deelnemers 15km bij het coördinatiepunt', 'Aanvang', 'aanvang', 10, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', null, null),
('1572e571-5f90-4cd9-b930-173d31df0124', '11:05u', 'Deelnemers 15km aanwezig startpunt (Kootwijk)', 'Aanwezig', 'aanwezig', 30, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', '52.18474064', '5.77074940'),
('16f07558-27cc-4e70-8d2f-4093d5e47009', '15:35u', 'START 2,5KM, Hervatting 6km, 10km en 15km', 'Start', 'start', 190, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', '52.22044762', '5.92889575'),
('2eb4673c-6bde-465a-b4e8-27425bc32d54', '15:15u', 'Verwachte aankomst 15, 10, 6 km lopers bij rustpunt (Berg & Bos - 15 min pauze)', 'Rustpunt', 'rustpunt', 180, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', null, null),
('425c461a-cee7-480c-a399-7d469cc7efbe', '12:50u', 'Deelnemers 10km aanwezig bij het startpunt (Halte Assel)', 'Aanwezig', 'aanwezig', 80, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', '52.20071362', '5.83602324'),
('460625b7-ec90-446d-a895-2ddaefb98335', '14:15u', 'START 6KM, Hervatting 10km en 15km', 'Start', 'start', 140, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', '52.21916438', '5.87255921'),
('4f201162-2524-4b03-b0ea-0fb5ccaf29c3', '15:00u', 'Vertrek deelnemers 2,5 km met de pendelbussen naar het startpunt 2,5km', 'Vertrek', 'vertrek', 160, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', null, null),
('5af9c5b7-9b73-44aa-9153-61e2277b9233', '17:00u – 18:00u', 'INHULDIGINGSFEEST', 'Feest', 'feest', 220, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', null, null),
('652b1d86-b062-4b63-bbf0-5294c71979d0', '12:45u', 'Verwachte aankomst 15 km lopers bij rustpunt (Halte Assel - 15 min pauze)', 'Rustpunt', 'rustpunt', 70, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', null, null),
('65dd7eb0-4de3-4d54-99b7-2f278a08ece4', '12:30u', 'Vertrek deelnemers 10km met de pendelbussen naar het startpunt 10km', 'Vertrek', 'vertrek', 60, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', null, null),
('6730216e-9a4c-4df5-aaca-119da8595eef', '10:45u', 'Vertrek pendelbussen naar startpunt 15km', 'Vertrek', 'vertrek', 20, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', null, null),
('7ae60214-dc2d-4b23-a4d1-0993fbb56e46', '14:00u', 'Deelnemers 6km aanwezig bij het startpunt (Hoog Soeren)', 'Aanwezig', 'aanwezig', 130, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', '52.21916438', '5.87255921'),
('91666bd3-32fd-4054-ab31-4d9f7532c0ce', '14:30u', 'Aanvang Deelnemers 2,5km bij het coördinatiepunt', 'Aanvang', 'aanvang', 150, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', null, null),
('9564d363-17d5-45b4-b868-45d165a82c72', '16:10u – 16:30u', 'FINISH', 'Finish', 'finish', 210, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', null, null),
('9ff53726-0e9a-48de-b4f9-35cfa7666756', '15:55u', 'Aankomst bij De Naald / START INHULDIGINGSLOOP', 'Aankomst', 'aankomst', 200, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', null, null),
('b9c046ea-de08-48e5-995e-9d4a98d76b6e', '14:00u', 'Verwachte aankomst 15, 10 km lopers bij rustpunt (Hoog Soeren - 15 min pauze)', 'Rustpunt', 'rustpunt', 120, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', null, null),
('bc30a0ea-0f22-443a-8124-0cd52e10a2b3', '11:15u', 'START 15KM', 'Start', 'start', 40, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', '52.18474064', '5.77074940'),
('cf1514db-2e06-45ab-891b-76736eb308d3', '13:00u', 'START 10KM, Hervatting 15km', 'Start', 'start', 90, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', '52.20071362', '5.83602324'),
('de0cd4c1-2fe0-4f13-9b94-00c03ce38527', '15:05u', 'Deelnemers 2,5km aanwezig bij het startpunt (Berg & Bos)', 'Aanwezig', 'aanwezig', 170, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', '52.22044762', '5.92889575'),
('fd80912e-f112-4bee-a2dd-570c4cf88c89', '13:15u', 'Aanvang Deelnemers 6km bij het coördinatiepunt', 'Aanvang', 'aanvang', 100, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', null, null),
('fec9b64c-4cdd-4fc1-a507-9c222fdeb958', '13:45u', 'Vertrek deelnemers 6 km met de pendelbussen naar het startpunt 6km', 'Vertrek', 'vertrek', 110, true, '2025-04-15 21:57:08.444242+00', '2025-04-15 21:57:08.444242+00', null, null)
ON CONFLICT (id) DO NOTHING;

-- Deel G: Social Embeds (van V1_31G)
CREATE TABLE IF NOT EXISTS social_embeds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    platform TEXT NOT NULL,
    embed_code TEXT NOT NULL,
    order_number INTEGER,
    visible BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

INSERT INTO social_embeds (id, platform, embed_code, order_number, visible, created_at, updated_at) VALUES
('5709a899-ee12-4883-8b57-3a4d0e8543a6', 'instagram', '<blockquote class="instagram-media" data-instgrm-permalink="https://www.instagram.com/p/DJ6farVIh7G/?utm_source=ig_embed&utm_campaign=loading" data-instgrm-version="14" style=" background:#FFF; border:0; border-radius:3px; box-shadow:0 0 1px 0 rgba(0,0,0,0.5),0 1px 10px 0 rgba(0,0,0,0.15); margin: 1px; max-width:540px; min-width:326px; padding:0; width:99.375%; width:-webkit-calc(100% - 2px); width:calc(100% - 2px);"><div style="padding:16px;"> <a href="https://www.instagram.com/p/DJ6farVIh7G/?utm_source=ig_embed&utm_campaign=loading" style=" background:#FFFFFF; line-height:0; padding:0 0; text-align:center; text-decoration:none; width:100%;" target="_blank"> <div style=" display: flex; flex-direction: row; align-items: center;"> <div style="background-color: #F4F4F4; border-radius: 50%; flex-grow: 0; height: 40px; margin-right: 14px; width: 40px;"></div> <div style="display: flex; flex-direction: column; flex-grow: 1; justify-content: center;"> <div style=" background-color: #F4F4F4; border-radius: 4px; flex-grow: 0; height: 14px; margin-bottom: 6px; width: 100px;"></div> <div style=" background-color: #F4F4F4; border-radius: 4px; flex-grow: 0; height: 14px; width: 60px;"></div></div></div><div style="padding: 19% 0;"></div> <div style="display:block; height:50px; margin:0 auto 12px; width:50px;"><svg width="50px" height="50px" viewBox="0 0 60 60" version="1.1" xmlns="https://www.w3.org/2000/svg" xmlns:xlink="https://www.w3.org/1999/xlink"><g stroke="none" stroke-width="1" fill="none" fill-rule="evenodd"><g transform="translate(-511.000000, -20.000000)" fill="#000000"><g><path d="M556.869,30.41 C554.814,30.41 553.148,32.076 553.148,34.131 C553.148,36.186 554.814,37.852 556.869,37.852 C558.924,37.852 560.59,36.186 560.59,34.131 C560.59,32.076 558.924,30.41 556.869,30.41 M541,60.657 C535.114,60.657 530.342,55.887 530.342,50 C530.342,44.114 535.114,39.342 541,39.342 C546.887,39.342 551.658,44.114 551.658,50 C551.658,55.887 546.887,60.657 541,60.657 M541,33.886 C532.1,33.886 524.886,41.1 524.886,50 C524.886,58.899 532.1,66.113 541,66.113 C549.9,66.113 557.115,58.899 557.115,50 C557.115,41.1 549.9,33.886 541,33.886 M565.378,62.101 C565.244,65.022 564.756,66.606 564.346,67.663 C563.803,69.06 563.154,70.057 562.106,71.106 C561.058,72.155 560.06,72.803 558.662,73.347 C557.607,73.757 556.021,74.244 553.102,74.378 C549.944,74.521 548.997,74.552 541,74.552 C533.003,74.552 532.056,74.521 528.898,74.378 C525.979,74.244 524.393,73.757 523.338,73.347 C521.94,72.803 520.942,72.155 519.894,71.106 C518.846,70.057 518.197,69.06 517.654,67.663 C517.244,66.606 516.755,65.022 516.623,62.101 C516.479,58.943 516.448,57.996 516.448,50 C516.448,42.003 516.479,41.056 516.623,37.899 C516.755,34.978 517.244,33.391 517.654,32.338 C518.197,30.938 518.846,29.942 519.894,28.894 C520.942,27.846 521.94,27.196 523.338,26.654 C524.393,26.244 525.979,25.756 528.898,25.623 C532.057,25.479 533.004,25.448 541,25.448 C548.997,25.448 549.943,25.479 553.102,25.623 C556.021,25.756 557.607,26.244 558.662,26.654 C560.06,27.196 561.058,27.846 562.106,28.894 C563.154,29.942 563.803,30.938 564.346,32.338 C564.756,33.391 565.244,34.978 565.378,37.899 C565.522,41.056 565.552,42.003 565.552,50 C565.552,57.996 565.522,58.943 565.378,62.101 M570.82,37.631 C570.674,34.438 570.167,32.258 569.425,30.349 C568.659,28.377 567.633,26.702 565.965,25.035 C564.297,23.368 562.623,22.342 560.652,21.575 C558.743,20.834 556.562,20.326 553.369,20.18 C550.169,20.033 549.148,20 541,20 C532.853,20 531.831,20.033 528.631,20.18 C525.438,20.326 523.257,20.834 521.349,21.575 C519.376,22.342 517.703,23.368 516.035,25.035 C514.368,26.702 513.342,28.377 512.574,30.349 C511.834,32.258 511.326,34.438 511.181,37.631 C511.035,40.831 511,41.851 511,50 C511,58.147 511.035,59.17 511.181,62.369 C511.326,65.562 511.834,67.743 512.574,69.651 C513.342,71.625 514.368,73.296 516.035,74.965 C517.703,76.634 519.376,77.658 521.349,78.425 C523.257,79.167 525.438,79.673 528.631,79.82 C531.831,79.965 532.853,80.001 541,80.001 C549.148,80.001 550.169,79.965 553.369,79.82 C556.562,79.673 558.743,79.167 560.652,78.425 C562.623,77.658 564.297,76.634 565.965,74.965 C567.633,73.296 568.659,71.625 569.425,69.651 C570.167,67.743 570.674,65.562 570.82,62.369 C570.966,59.17 571,58.147 571,50 C571,41.851 570.966,40.831 570.82,37.631"></path></g></g></g></svg></div><div style="padding-top: 8px;"> <div style=" color:#3897f0; font-family:Arial,sans-serif; font-size:14px; font-style:normal; font-weight:550; line-height:18px;">Dit bericht op Instagram bekijken</div></div><div style="padding: 12.5% 0;"></div> <div style="display: flex; flex-direction: row; margin-bottom: 14px; align-items: center;"><div> <div style="background-color: #F4F4F4; border-radius: 50%; height: 12.5px; width: 12.5px; transform: translateX(0px) translateY(7px);"></div> <div style="background-color: #F4F4F4; height: 12.5px; transform: rotate(-45deg) translateX(3px) translateY(1px); width: 12.5px; flex-grow: 0; margin-right: 14px; margin-left: 2px;"></div> <div style="background-color: #F4F4F4; border-radius: 50%; height: 12.5px; width: 12.5px; transform: translateX(9px) translateY(-18px);"></div></div><div style="margin-left: 8px;"> <div style=" background-color: #F4F4F4; flex-grow: 0; height: 20px; width: 20px;"></div> <div style=" width: 0; height: 0; border-top: 2px solid transparent; border-left: 6px solid #f4f4f4; border-bottom: 2px solid transparent; transform: translateX(16px) translateY(-4px) rotate(30deg)"></div></div><div style="margin-left: auto;"> <div style=" width: 0px; border-top: 8px solid #F4F4F4; border-right: 8px solid transparent; transform: translateY(16px);"></div> <div style=" background-color: #F4F4F4; flex-grow: 0; height: 12px; width: 16px; transform: translateY(-4px);"></div> <div style=" width: 0; height: 0; border-top: 8px solid #F4F4F4; border-left: 8px solid transparent; transform: translateY(-4px) translateX(8px);"></div></div></div> <div style="display: flex; flex-direction: column; flex-grow: 1; justify-content: center; margin-bottom: 24px;"> <div style=" background-color: #F4F4F4; border-radius: 4px; flex-grow: 0; height: 14px; margin-bottom: 6px; width: 224px;"></div> <div style=" background-color: #F4F4F4; border-radius: 4px; flex-grow: 0; height: 14px; width: 144px;"></div></div></a><p style=" color:#c9c8cd; font-family:Arial,sans-serif; font-size:14px; line-height:17px; margin-bottom:0; margin-top:8px; overflow:hidden; padding:8px 0 7px; text-align:center; text-overflow:ellipsis; white-space:nowrap;"><a href="https://www.instagram.com/p/DJ6farVIh7G/?utm_source=ig_embed&utm_campaign=loading" style=" color:#c9c8cd; font-family:Arial,sans-serif; font-size:14px; font-style:normal; font-weight:normal; line-height:17px; text-decoration:none;" target="_blank">Een bericht gedeeld door Koninklijke Loop (@koninklijkeloop)</a></p></div></blockquote>
<script async src="//www.instagram.com/embed.js"></script>', 1, true, '2024-12-22 19:30:53.366424+00', '2024-12-22 19:30:53.366424+00'),
('ee8d1152-2fd4-464d-82a8-76c02ad56ed9', 'facebook', '<iframe src="https://www.facebook.com/plugins/post.php?href=https%3A%2F%2Fwww.facebook.com%2Fpermalink.php%3Fstory_fbid%3Dpfbid02XNU75Y2gMxWhvsVQar7oaM98GvMLLryXQVMTjxnBkEg6e6imJ8ecgoEF9SrTVJDpl%26id%3D61556315443279&show_text=true&width=500" width="500" height="737" style="border:none;overflow:hidden" scrolling="no" frameborder="0" allowfullscreen="true" allow="autoplay; clipboard-write; encrypted-media; picture-in-picture; web-share"></iframe>', 2, true, '2024-12-22 19:30:53.366424+00', '2024-12-22 19:30:53.366424+00')
ON CONFLICT (id) DO NOTHING;

-- Deel H: Social Links (van V1_31H)
CREATE TABLE IF NOT EXISTS social_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    platform TEXT NOT NULL,
    url TEXT NOT NULL,
    bg_color_class TEXT,
    icon_color_class TEXT,
    order_number INTEGER,
    visible BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

INSERT INTO social_links (id, platform, url, bg_color_class, icon_color_class, order_number, visible, created_at, updated_at) VALUES
('1de3ed0c-9bf3-4924-8d76-72045cc1c0ec', 'instagram', 'https://www.instagram.com/koninklijkeloop', null, null, 2, true, '2024-12-22 19:49:18.04807+00', '2024-12-22 19:49:18.04807+00'),
('29988672-0d98-4083-9908-18f6bbe34f3f', 'linkedin', 'https://www.linkedin.com/company/koninklijkeloop', null, null, 4, true, '2024-12-22 19:49:18.04807+00', '2024-12-22 19:49:18.04807+00'),
('dc917a65-6bb4-45cc-a1dc-890b1cdf1f5b', 'facebook', 'https://www.facebook.com/koninklijkeloop', null, null, 1, true, '2024-12-22 19:49:18.04807+00', '2024-12-22 19:49:18.04807+00'),
('f713d59e-e8a8-40d2-bb11-fc64f89b9ae9', 'youtube', 'https://www.youtube.com/@koninklijkeloop', null, null, 3, true, '2024-12-22 19:49:18.04807+00', '2024-12-22 19:49:18.04807+00')
ON CONFLICT (id) DO NOTHING;

-- Deel I: Under Construction (van V1_31I)
CREATE TABLE IF NOT EXISTS under_construction (
    id SERIAL PRIMARY KEY,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    title TEXT,
    message TEXT,
    footer_text TEXT,
    logo_url TEXT,
    expected_date TIMESTAMP WITH TIME ZONE,
    social_links JSONB,
    progress_percentage INTEGER,
    contact_email TEXT,
    newsletter_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

INSERT INTO under_construction (id, is_active, title, message, footer_text, logo_url, expected_date, social_links, progress_percentage, contact_email, newsletter_enabled, created_at, updated_at) VALUES
(1, false, 'Website in onderhoud', 'We stomen ons klaar voor De Koninklijke Loop 2026, op dit moment is de website helaas niet bereikbaar', 'Bedankt voor uw geduld!', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png', '2026-01-31 18:00:00+00', '[{"url": "https://twitter.com/koninklijkeloop", "platform": "Twitter"}, {"url": "https://instagram.com/koninklijkeloop", "platform": "Instagram"}, {"url": "https://www.youtube.com/@DeKoninklijkeLoop", "platform": "YouTube"}]', 85, 'info@koninklijkeloop.nl', false, '2025-09-26 17:37:22.197854+00', '2025-10-09 21:00:29.391392+00')
ON CONFLICT (id) DO NOTHING;

-- Deel J: Partners (van V1_32)
CREATE TABLE IF NOT EXISTS partners (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    logo TEXT,
    website TEXT,
    tier TEXT,
    since DATE,
    visible BOOLEAN NOT NULL DEFAULT TRUE,
    order_number INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Deel K: Radio Recordings (van V1_32)
CREATE TABLE IF NOT EXISTS radio_recordings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    description TEXT,
    date TEXT,
    audio_url TEXT,
    thumbnail_url TEXT,
    visible BOOLEAN NOT NULL DEFAULT TRUE,
    order_number INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Insert partners data (idempotent)
INSERT INTO "public"."partners" ("id", "name", "description", "logo", "website", "tier", "since", "visible", "order_number", "created_at", "updated_at") VALUES
('26f11d04-c0da-4755-b2e1-5fcb5e9887d4', 'Accress', 'Accres beheert in Apeldoorn ruim 60 locaties, waaronder sporthallen, wijkcentra, zwembaden, kinderboerderijen en een stadspark.

Zij helpen ond bij het halen van ons doel!', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1744388421/accres_logo_ochsmg.jpg', 'https://www.accres.nl/', 'bronze', '2025-04-01', 'true', '5', '2025-04-11 16:21:40+00', '2025-04-11 16:45:10.227784+00'),
('510d8e4d-6a7d-4ab3-b311-17b32df7f01f', 'Apeldoorn', 'Apeldoorn ondersteunt ons in ons doel.', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1734895194/nw1qxouzupzsshckzkab.png', 'https://www.apeldoorn.nl/', 'bronze', '2024-12-22', 'true', '0', '2024-12-22 19:19:55.414207+00', '2025-02-26 21:10:59.040907+00'),
('5e8b8390-e637-4c6d-80d8-25ff4155d5a9', 'Sheeren Loo', 'Samen met bewoners van SheerenLoo wordt deze loop georganiseerd.', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1734894570/mtvucaouruenat2cllsi.png', 'https://www.sheerenloo.nl/', 'silver', '2024-12-22', 'true', '0', '2024-12-22 19:09:31.099813+00', '2025-01-07 14:11:39.051238+00'),
('9c311d0f-4db7-4da9-9d87-a836e48078cd', 'Liliane Fonds', 'Samen maken we ons sterk', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1734893709/qsygajx2tdxxbqbfyurr.png', 'https://www.lilianefonds.nl/', 'bronze', '2024-12-22', 'true', '0', '2024-12-22 18:55:09.823207+00', '2025-01-08 19:43:58.464461+00'),
('eda49448-3db9-4db5-9d01-a7331879f5b5', 'De Grote Kerk', 'De grote kerk ondersteund ons al vanaf het begin. Hier is het allemaal begonnen.', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1734895146/ri4vclttn4nn2wh53wj0.jpg', 'https://www.grotekerkapeldoorn.nl/', 'gold', '2024-12-22', 'true', '0', '2024-12-22 19:19:07.18416+00', '2025-04-11 16:46:10.191766+00')
ON CONFLICT (id) DO NOTHING;

-- Insert radio_recordings data (idempotent)
INSERT INTO "public"."radio_recordings" ("id", "title", "description", "date", "audio_url", "thumbnail_url", "visible", "order_number", "created_at", "updated_at") VALUES
('a6e73425-6af2-4bbd-84f8-67b16d195f99', 'De koninklijke Loop 2025 uitzending!', 'Luister naar het live radioverslag van De Koninklijke Loop 2025, uitgezonden op RTV Apeldoorn, met interviews met de organisatie!

', '14 mei 2025', 'https://res.cloudinary.com/dgfuv7wif/video/upload/v1747733438/DKLRTV2025_dpdydc.wav', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png', 'true', '1', '2025-05-20 09:32:01+00', '2025-05-20 09:33:27.540519+00'),
('c2d84ca9-97a9-4990-9ee4-1fe1718a8c5b', 'Radioverslag Koninklijke Loop 2024 (RTV Apeldoorn)', 'Luister naar het live radioverslag van De Koninklijke Loop 2024, uitgezonden op RTV Apeldoorn, met interviews en sfeerimpressies.', '15 mei 2024', 'https://res.cloudinary.com/dgfuv7wif/video/upload/v1714042357/matinee_1_nbm0ph.wav', 'https://res.cloudinary.com/dgfuv7wif/image/upload/v1733267882/664b8c1e593a1e81556b4238_0760849fb8_yn6vdm.png', 'true', '2', '2025-04-08 19:50:35.125804+00', '2025-05-20 09:32:30.549188+00')
ON CONFLICT (id) DO NOTHING;