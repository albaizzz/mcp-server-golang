IF OBJECT_ID(N'dbo.kodepos', N'U') IS NOT NULL
  DROP TABLE dbo.kodepos;
GO

CREATE TABLE dbo.kodepos (
    id            BIGINT IDENTITY(1,1) PRIMARY KEY,
    kodepos       CHAR(5)     NOT NULL,
    kelurahan     NVARCHAR(100) NOT NULL,
    kecamatan     NVARCHAR(100) NOT NULL,
    kota_kab      NVARCHAR(120) NOT NULL,
    provinsi      NVARCHAR(120) NOT NULL,
    latitude      DECIMAL(9,6)  NULL,
    longitude     DECIMAL(9,6)  NULL
);