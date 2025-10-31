
--
-- PostgreSQL database dump
--

-- Dumped from database version 14.6 (Ubuntu 14.6-0ubuntu0.22.04.1)
-- Dumped by pg_dump version 14.6 (Ubuntu 14.6-0ubuntu0.22.04.1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: goods; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.goods (
    id integer NOT NULL,
    good character varying(255),
    price integer,
    count integer
);


--
-- Name: goods_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.goods_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: goods_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.goods_id_seq OWNED BY public.goods.id;


--
-- Name: orders; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.orders (
    good_id integer,
    status_id integer
);


--
-- Name: goods id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.goods ALTER COLUMN id SET DEFAULT nextval('public.goods_id_seq'::regclass);


--
-- Data for Name: goods; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.goods (id, good, price, count) FROM stdin;
1	core i5 4050	950	123
2	core i3 3010	350	432
3	celeron 5150	150	12
4	ryzen 5 1600	153	30
5	ryzen 5 1700	850	650
6	ryzen 7 5800	1230	250
7	ryzen 7 5600X	900	0
8	ryzen 3 2200G	200	0
9	ryzen 3 3500G	200	3
10	ryzen 9 5900	1700	5
\.


--
-- Data for Name: orders; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.orders (good_id, status_id) FROM stdin;
6	0
8	0
9	7
3	7
1	0
2	7
\.


--
-- Name: goods_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.goods_id_seq', 10, true);


--
-- Name: goods goods_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.goods
    ADD CONSTRAINT goods_pkey PRIMARY KEY (id);


--
-- Name: orders orders_good_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_good_id_fkey FOREIGN KEY (good_id) REFERENCES public.goods(id);


--
-- PostgreSQL database dump complete
--

