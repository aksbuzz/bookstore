DROP TABLE IF EXISTS "book";
DROP SEQUENCE IF EXISTS book_id_seq;
CREATE SEQUENCE book_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."book" (
    "id" integer DEFAULT nextval('book_id_seq') NOT NULL,
    "author" character varying(255) NOT NULL,
    "category" character varying(255) NOT NULL,
    "cover" character varying(255) DEFAULT '',
    "name" character varying(255) NOT NULL,
    "price" double precision NOT NULL,
    "rating" double precision DEFAULT '0' NOT NULL,
    CONSTRAINT "book_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

CREATE INDEX "book_category" ON "public"."book" USING btree ("category");


DROP TABLE IF EXISTS "cart";
DROP SEQUENCE IF EXISTS cart_id_seq;
CREATE SEQUENCE cart_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."cart" (
    "id" integer DEFAULT nextval('cart_id_seq') NOT NULL,
    "book_id" integer NOT NULL,
    "quantity" smallint NOT NULL,
    "price" double precision NOT NULL,
    CONSTRAINT "cart_pkey" PRIMARY KEY ("id")
) WITH (oids = false);


DROP TABLE IF EXISTS "order";
DROP SEQUENCE IF EXISTS order_id_seq;
CREATE SEQUENCE order_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."order" (
    "id" integer DEFAULT nextval('order_id_seq') NOT NULL,
    "date" timestamp NOT NULL,
    "total" double precision NOT NULL,
    CONSTRAINT "order_pkey" PRIMARY KEY ("id")
) WITH (oids = false);


DROP TABLE IF EXISTS "order_item";
DROP SEQUENCE IF EXISTS order_item_id_seq;
CREATE SEQUENCE order_item_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."order_item" (
    "id" integer DEFAULT nextval('order_item_id_seq') NOT NULL,
    "order_id" integer NOT NULL,
    "book_id" integer NOT NULL,
    "quantity" smallint NOT NULL,
    "price" double precision NOT NULL,
    CONSTRAINT "order_item_pkey" PRIMARY KEY ("id")
) WITH (oids = false);


ALTER TABLE ONLY "public"."cart" ADD CONSTRAINT "cart_book_id_fkey" FOREIGN KEY (book_id) REFERENCES book(id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."order_item" ADD CONSTRAINT "order_item_book_id_fkey" FOREIGN KEY (book_id) REFERENCES book(id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."order_item" ADD CONSTRAINT "order_item_order_id_fkey" FOREIGN KEY (order_id) REFERENCES "order"(id) NOT DEFERRABLE;

