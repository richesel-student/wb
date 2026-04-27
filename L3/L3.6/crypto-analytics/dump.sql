--
-- PostgreSQL database dump
--

\restrict yolvqyvnBXJzwPoVq6UOJcIa051EgpuABc0CbvSSdz85Gr8DQ41pCa9H1zFPVun

-- Dumped from database version 18.3 (Homebrew)
-- Dumped by pg_dump version 18.3 (Homebrew)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
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
-- Name: transactions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.transactions (
    id integer NOT NULL,
    asset_id integer,
    type text,
    amount_usd numeric,
    price numeric,
    quantity numeric,
    fee numeric,
    "timestamp" timestamp without time zone DEFAULT now(),
    asset_name text DEFAULT 'BTC'::text
);


--
-- Name: transactions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.transactions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: transactions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.transactions_id_seq OWNED BY public.transactions.id;


--
-- Name: transactions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.transactions ALTER COLUMN id SET DEFAULT nextval('public.transactions_id_seq'::regclass);


--
-- Data for Name: transactions; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.transactions (id, asset_id, type, amount_usd, price, quantity, fee, "timestamp", asset_name) FROM stdin;
12	\N	buy	100	76956.5	0.0012994353953207332	0.1	2026-04-27 23:40:41.351602	BTC
13	\N	buy	100	76956.5	0.0012994353953207332	0.1	2026-04-27 23:40:50.539858	ETH
14	\N	buy	100	76959.6	0.0012993830529264704	0.1	2026-04-27 23:40:52.128453	ETH
15	\N	sell	100	76967.7	0.0012992463072171833	0.1	2026-04-27 23:40:56.770758	ETH
9	\N	buy	2000	76691.7	0.02607844134371777	2	2026-04-27 00:25:22.073514	ETH
10	\N	buy	2000	76690.7	0.026078781390703177	2	2026-04-26 00:25:22.073764	SOL
11	\N	buy	2000	76836.6	0.026029262096448826	2	2026-04-25 00:25:22.073997	BTC
\.


--
-- Name: transactions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.transactions_id_seq', 15, true);


--
-- Name: transactions transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

\unrestrict yolvqyvnBXJzwPoVq6UOJcIa051EgpuABc0CbvSSdz85Gr8DQ41pCa9H1zFPVun

