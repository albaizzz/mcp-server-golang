-- Pencarian by kodepos persis
CREATE INDEX IX_kodepos_kodepos ON dbo.kodepos (kodepos);

-- Pencarian teks (prefix) pada beberapa kolom
CREATE INDEX IX_kodepos_kota_kec_kel
ON dbo.kodepos (kota_kab, kecamatan, kelurahan);

-- Jika dataset besar, pertimbangkan full-text (opsional)
-- CREATE FULLTEXT CATALOG ft AS DEFAULT;
-- CREATE FULLTEXT INDEX ON dbo.kodepos (kelurahan LANGUAGE 1057, kecamatan LANGUAGE 1057, kota_kab LANGUAGE 1057, provinsi LANGUAGE 1057)
-- KEY INDEX PK__kodepos__id;