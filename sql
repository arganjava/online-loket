--
-- PostgreSQL database dump
--

-- Dumped from database version 12.3
-- Dumped by pg_dump version 12.3

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
-- Name: event; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.event (
    id character varying(100) NOT NULL,
    event_name character varying(30),
    description character varying(200),
    schedule_begin date,
    schedule_end date,
    location_id character varying(100)
);


ALTER TABLE public.event OWNER TO postgres;

--
-- Name: event_ticket; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.event_ticket (
    id character varying(100) NOT NULL,
    event_id character varying(100),
    ticket_type character varying(5),
    quantity integer,
    price numeric
);


ALTER TABLE public.event_ticket OWNER TO postgres;

--
-- Name: location; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.location (
    country character varying(30),
    city_name character varying(50),
    id character varying(100) NOT NULL,
    address character varying(200),
    village character varying(50)
);


ALTER TABLE public.location OWNER TO postgres;

--
-- Name: transaction; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.transaction (
    event_ticket_id character varying(100),
    id character varying(100) NOT NULL,
    customer_name character varying(100),
    customer_phone character varying(20),
    customer_email character varying(100),
    order_quantity integer,
    transaction_time timestamp without time zone
);


ALTER TABLE public.transaction OWNER TO postgres;

--
-- Data for Name: event; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.event (id, event_name, description, schedule_begin, schedule_end, location_id) FROM stdin;
30c71f1a-da5b-43df-a8fc-091c8e4452d7	Badut Kota	Desc	2020-12-30	2021-01-30	1ad5ae0e-9e49-4025-90aa-295e1a4bd886
\.


--
-- Data for Name: event_ticket; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.event_ticket (id, event_id, ticket_type, quantity, price) FROM stdin;
9fff2621-ecf7-42f4-9125-524a90b868a9	30c71f1a-da5b-43df-a8fc-091c8e4452d7	CHILD	10	1000
ae5733f8-3d20-40d5-b6b2-caa56e2d36c8	30c71f1a-da5b-43df-a8fc-091c8e4452d7	ADULT	8	1000
\.


--
-- Data for Name: location; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.location (country, city_name, id, address, village) FROM stdin;
Indonesia	Jakarta	69707ecf-0dfa-49e4-84c8-36755003ff2a	Address	village1
Indonesia	Bandung	1ad5ae0e-9e49-4025-90aa-295e1a4bd886	Address	Ujung Berung
Indonesia	Bekasi	c2366d15-2b4e-471c-8785-f5ffa5f06184	Address	Test
\.


--
-- Data for Name: transaction; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.transaction (event_ticket_id, id, customer_name, customer_phone, customer_email, order_quantity, transaction_time) FROM stdin;
ae5733f8-3d20-40d5-b6b2-caa56e2d36c8	4f2836ed-c9ff-471b-b925-8992c1545d5e	Argan Megariansyah	0123232	arganjava@gmail.com	2	2020-12-13 15:48:38.309315
\.


--
-- Name: event event_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event
    ADD CONSTRAINT event_pk PRIMARY KEY (id);


--
-- Name: event_ticket event_ticket_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event_ticket
    ADD CONSTRAINT event_ticket_pk PRIMARY KEY (id);


--
-- Name: location location_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.location
    ADD CONSTRAINT location_pk PRIMARY KEY (id);


--
-- Name: transaction transaction_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transaction
    ADD CONSTRAINT transaction_pk PRIMARY KEY (id);


--
-- Name: event event_location_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event
    ADD CONSTRAINT event_location_id_fk FOREIGN KEY (location_id) REFERENCES public.location(id);


--
-- Name: event_ticket event_ticket_event_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event_ticket
    ADD CONSTRAINT event_ticket_event_id_fk FOREIGN KEY (event_id) REFERENCES public.event(id);


--
-- Name: transaction transaction_event_ticket_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transaction
    ADD CONSTRAINT transaction_event_ticket_id_fk FOREIGN KEY (event_ticket_id) REFERENCES public.event_ticket(id);


--
-- PostgreSQL database dump complete
--

