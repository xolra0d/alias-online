--
-- PostgreSQL database dump
--

\restrict hjMldxUvJNq16ROR1Ywn56I5d8qbxk2RKLJkXnC3C5ioF0ExZeqFkUl2JPZPXUk

-- Dumped from database version 18.2
-- Dumped by pg_dump version 18.2

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
-- Name: room_participants; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.room_participants (
    room_id text NOT NULL,
    user_id uuid NOT NULL,
    words_tried integer DEFAULT 0,
    words_guessed integer DEFAULT 0,
    turn_order integer NOT NULL
);


ALTER TABLE public.room_participants OWNER TO postgres;

--
-- Name: room_participants_turn_order_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.room_participants_turn_order_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.room_participants_turn_order_seq OWNER TO postgres;

--
-- Name: room_participants_turn_order_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.room_participants_turn_order_seq OWNED BY public.room_participants.turn_order;


--
-- Name: rooms; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.rooms (
    id text NOT NULL,
    admin uuid NOT NULL,
    seed integer NOT NULL,
    current_word_index integer DEFAULT 0 NOT NULL,
    current_player_id uuid NOT NULL,
    game_state integer DEFAULT 0,
    language text NOT NULL,
    rude_words boolean NOT NULL,
    additional_vocabulary text[] DEFAULT '{}'::text[] NOT NULL,
    clock integer NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone
);


ALTER TABLE public.rooms OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id uuid NOT NULL,
    name text NOT NULL,
    secret_hash text NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    last_seen timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: vocabularies; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.vocabularies (
    language text NOT NULL,
    primary_words text[] NOT NULL,
    rude_words text[] DEFAULT '{}'::text[] NOT NULL,
    available boolean DEFAULT true NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone
);


ALTER TABLE public.vocabularies OWNER TO postgres;

--
-- Name: room_participants turn_order; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.room_participants ALTER COLUMN turn_order SET DEFAULT nextval('public.room_participants_turn_order_seq'::regclass);


--
-- Data for Name: room_participants; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.room_participants (room_id, user_id, words_tried, words_guessed, turn_order) FROM stdin;
DXG2NKOK	019cec7d-3de9-74fe-8dc8-f881424593ad	0	0	0
DXG2NKOK	019cf293-62fb-73ac-b73b-cc883361244f	0	0	1
YCEGHS6F	019cec7d-3de9-74fe-8dc8-f881424593ad	34	1	0
YCEGHS6F	019cf293-62fb-73ac-b73b-cc883361244f	0	0	1
GWFMNN3T	019cf293-62fb-73ac-b73b-cc883361244f	17	0	0
GWFMNN3T	019cec7d-3de9-74fe-8dc8-f881424593ad	0	0	1
\.


--
-- Data for Name: rooms; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.rooms (id, admin, seed, current_word_index, current_player_id, game_state, language, rude_words, additional_vocabulary, clock, created_at, updated_at) FROM stdin;
DXG2NKOK	019cec7d-3de9-74fe-8dc8-f881424593ad	301328077	0	019cf293-62fb-73ac-b73b-cc883361244f	0	English	f	{}	60	2026-03-16 14:29:33.821711	\N
YCEGHS6F	019cec7d-3de9-74fe-8dc8-f881424593ad	883850484	33	019cec7d-3de9-74fe-8dc8-f881424593ad	1	English	f	{}	60	2026-03-16 14:35:27.14494	\N
GWFMNN3T	019cf293-62fb-73ac-b73b-cc883361244f	1519597312	17	019cec7d-3de9-74fe-8dc8-f881424593ad	0	English	f	{latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif,latif}	60	2026-03-16 14:36:37.480007	\N
3AXJBOJA	019cec7d-3de9-74fe-8dc8-f881424593ad	163699406	0	019cec7d-3de9-74fe-8dc8-f881424593ad	0	English	f	{}	60	2026-03-16 14:51:02.755633	\N
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, name, secret_hash, created_at, last_seen) FROM stdin;
019cd727-d673-78f5-970b-691e94c825df	WobblyMuffin12	e57b5fcdf5f7d542f0555632e5676df427102913320a30621a8d8e954d284343	2026-03-10 11:50:38.222165	2026-03-10 11:50:38.222165
019ce73e-b82c-72b1-8767-47638142c9b5	SpicyWaffle91	95bdec0caca34bc6fc5c75afe265fb2b272bf7077c35f73750f566783467eff1	2026-03-13 14:49:33.247658	2026-03-13 14:49:33.247658
019ce7cd-307d-754a-8910-34938a25f165	ChaoticWizard27	e00898dedb66ad2b7f19ae04d9003eeeac137ccb0c3f97fad6e1980914694ca4	2026-03-13 17:25:10.142807	2026-03-13 17:25:10.142807
019cec7d-3de9-74fe-8dc8-f881424593ad	ChaoticBiscuit31	b40e3d3e8a2156379947ed0da783d8026e5be60fe21b3d29ea92ecfcdac76d8e	2026-03-14 15:15:56.802834	2026-03-14 15:15:56.802834
019cf293-62fb-73ac-b73b-cc883361244f	GrumpyNoodle48	fa1a3236cd29e111236dba77e32ff673abe7977369ebdd9cd37e8f5a71d37fc3	2026-03-15 19:37:51.355766	2026-03-15 19:37:51.355766
\.


--
-- Data for Name: vocabularies; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.vocabularies (language, primary_words, rude_words, available, created_at, updated_at) FROM stdin;
English	{pear,meteor,t-shirt,hook,thick,climbing,trap,boss,seal,cafe,broccoli,keyboard,failure,fast,chariot,ridge,boat,chipmunk,lava,egg,sled,orange,power,block,tiger,flint,open,cup,yew,origami,iceberg,log,crib,tusk,tunnel,jam,gnome,cone,sole,ruby,notebook,exit,rocket,wasp,den,binoculars,witch,blueberry,rose,car,salamander,bittern,salad,thorn,dune,scroll,monk,rough,crimson,cucumber,factory,freckle,canopy,empire,campfire,gorge,brine,dress,orca,supermarket,adventure,dam,bath,pebble,well,bus,jump,airport,auditorium,heavy,wrestling,quill,hurricane,bison,algorithm,hamster,bed,drill,dream,canoe,dwarf,pit,road,cave,fig,hotel,build,greenhouse,aquarium,battleship,rust,secret,lemon,clam,mule,spore,create,cement,currency,goat,barrel,pond,igloo,crocodile,riddle,bruise,soap,rabbit,mint,future,branch,soup,jewel,slide,citrus,cat,rugby,pickaxe,risk,ripple,lens,card,fence,badminton,chalk,milk,corner,big,cloud,boomerang,grid,sting,rock,axe,tulip,cocoon,panda,newt,carry,station,peace,lock,buttercup,web,depot,sphere,taste,pepper,market,cellar,prawn,drop,mill,plate,scooter,oasis,azure,funnel,drain,bridge,measure,belt,watermelon,pancake,balcony,ceiling,count,corn,tennis,valley,knowledge,envelope,passage,poppy,seed,watchtower,compare,dolphin,atlas,hut,cicada,reef,knife,chive,airplane,agree,whale,late,port,anvil,idea,soft,level,apple,spray,bounty,website,compass,field,lynx,ruin,skirt,hard,beach,chocolate,filter,boggle,slab,harbor,pixel,wide,dance,slice,crowd,wing,giraffe,panther,jeans,tide,bookmark,basketball,camera,paw,bombard,cage,willow,billboard,brick,catch,restaurant,signal,quarry,gloves,ocean,mix,loop,arrowhead,key,analyze,dye,treadmill,stump,horizon,beaver,coat,stripe,mask,serve,smooth,cedar,sail,wisdom,donkey,north,cobweb,spinach,circus,spider,hero,bonfire,dawn,narwhal,beard,mitten,vault,raven,napkin,captain,moose,hardware,dark,brush,monkey,crayon,basement,remote,steak,ditch,paint,boar,complex,cabin,artillery,shelf,fly,shiny,plume,comet,desert,bicycle,trough,grape,bloom,snare,wand,crater,daisy,flour,cow,herb,clay,pie,lift,garden,lime,clock,glass,grain,skating,rice,hose,tray,anchor,cypress,curtain,pen,horn,mist,dog,train,freedom,doughnut,column,download,flesh,gust,owl,marsh,beehive,backyard,leaf,candle,boil,cushion,perch,sponge,torch,kite,cobra,shadow,thunder,jellyfish,arch,table,promise,sandwich,stool,golf,lighthouse,past,carpet,wade,submarine,calendar,shout,bone,wheel,mouse,moon,bagel,blackboard,plain,brittle,swan,crate,crop,ring,slug,waffle,hoodie,elephant,mystery,coral,talk,mountain,bottle,stockpile,lobster,ancient,sand,pasta,volcano,hat,bake,pineapple,goal,layer,dart,slope,shark,modern,cinder,stone,shampoo,oar,papaya,condor,stream,castle,chair,petal,truth,pouch,chimney,complain,zeppelin,explosion,grass,eclipse,cactus,hedge,wave,crossroads,bulletin,satellite,skiing,braid,finger,bedrock,fog,internet,flower,pig,shirt,curry,lily,cycling,dragon,python,ship,frost,hill,football,"ice cream",parrot,stallion,sushi,application,parachute,carrot,flamingo,cut,mop,guide,forest,cap,eel,mole,mongoose,staircase,motorcycle,wardrobe,rain,narrow,birch,wind,haystack,spring,koala,honey,butter,storm,caravan,cliff,viper,jungle,baby,sofa,rapids,surge,fang,bacon,chip,floor,skeleton,gym,hyena,bean,jacket,wilderness,pizza,creek,camel,orbit,chart,cougar,crown,yacht,surfing,plan,database,meadow,deer,cook,puma,scale,otter,swimming,journal,crest,hump,pumpkin,jelly,nymph,tram,bow,bookcase,knot,claw,chicken,gulf,harvest,fleet,attic,twig,bubble,browser,baseball,flag,ivy,laptop,blueprint,autumn,blowfish,scaffold,brutus,catalyst,glue,raft,waterfall,hockey,yarn,whistle,snail,ferry,glow,push,pillow,pirate,hammer,chance,fork,jaguar,spoon,drum,running,memory,sky,contest,cheek,pour,insect,cart,queen,plum,mirror,scorpion,turtle,mushroom,kitten,granite,vein,bolt,truck,melon,duck,canyon,thin,wreck,run,draw,firework,colorful,slippers,sprout,galaxy,engine,pencil,username,shorts,slow,spaceship,tall,chancellor,flame,cyclone,taxi,cocoa,fur,whirlpool,champion,palm,aurora,avocado,break,chorus,archive,chrome,email,swift,spike,quiver,walk,update,bowl,"beach ball",turkey,swamp,nest,booth,root,lotus,worm,croc,fudge,dove,broom,kingdom,lawn,"polar bear",cascade,sun,shrimp,bush,blush,seagull,snake,wizard,penguin,ask,tower,knight,iris,pull,write,tree,speaker,mug,swim,cobblestone,beacon,small,climb,farm,grove,zone,paddle,cloth,short,vacuum,prism,city,chain,bread,bay,raccoon,office,ledge,map,stem,basin,eraser,island,peach,blizzard,cricket,banana,sweater,vapor,region,stadium,soil,shovel,"coral reef",rope,ramp,flood,stick,minnow,light,charcoal,magnet,scarf,museum,guard,nut,river,drawer,toothbrush,helicopter,croak,bench,dice,probe,pigeon,bamboo,hive,hollow,toad,capsule,buffalo,success,library,pharmacy,cube,kayak,lantern,lagoon,spawn,throw,magic,toothpaste,obsidian,copper,torrent,cinema,splash,lake,bucket,fire,cleaner,coyote,cheese,backpack,star,border,blacksmith,tsunami,lightning,tundra,vulture,ink,hippo,ruler,basket,porch,villain,yak,bear,bakery,noodle,ember,screen,armor,password,garlic,cookie,biome,walrus,helmet,brook,fold,pudding,prairie,candy,school,fall,bright,old,puzzle,pine,crab,cake,taco,whisper,hotdog,dome,close,bat,plank,bramble,berry,badge,bulb,bank,sword,explain,courage,wolf,mango,pale,arrow,chasm,quake,robin,sheep,zoom,cape,robot,crawl,shoes,swallow,pelican,early,venom,brain,subway,cheetah,flask,fern,closet,microscope,answer,long,phone,server,towel,cluster,scissors,wheat,pasture,shell,fist,ranger,convoy,barn,charger,escalator,hawk,rift,fin,fossil,starfish,beam,chef,jaw,chrysalis,bell,clown,battery,boxing,frisbee,bellow,skull,young,eagle,crossbow,network,software,hay,village,hub,crystals,"artificial intelligence",frog,crew,treasure,upload,leopard,rye,sycamore,couch,zebra,blade,volleyball,boulder,moss,vine,socks,tribe,justice,wig,lamp,octopus,mud,trident,fox,kangaroo,step,burger,barricade,simple,chamber,park,snow,llama,burrow,balloon,gauge,feather,spark,gem,argue,blanket,spiral,glacier,dive,luck,snowboarding,bunny,scarecrow,"zebra fish",hound,saddle,olive,lion,trail,fish,horse,seaweed,baboon,"wolf cub",squirrel,door,blossom,gate,arena,yogurt,mangrove,caterpillar,skunk,cherry,hospital,axle,cockroach,sandals,read,boots,theater,present,blast,fry}	{douchebag,jerkface,goof,dingbat,geek,freeloader,wannabe,meathead,airhead,moron,cheater,dummy,leech,dipstick,screwball,dumbass,backstabber,brat,screw-up,tryhard,"loud jerk",jackass,jackwagon,snob,smartass,dimwit,scumbag,hater,weirdo,suck-up,prick,"drama queen",bigmouth,bastard,loudmouth,hypocrite,parasite,bully,liar,loser,fraud,knucklehead,rat,snake,bitch,shithead,nerd,troublemaker,idiot,jerk,pig,douche,asshole,clownboy,"keyboard warrior",blockhead,slacker,pain-in-the-neck,bonehead,troll,creep,fool,crybully,nitwit,brown-noser,clown,pain-in-the-ass,joker,halfwit,coward,numbskull,badass,dumbshit,birdbrain,"son of a bitch",fake,crybaby,show-off,sleazebag,poser}	t	2026-03-14 10:20:09.753158	\N
Own vocabulary	{}	{}	t	2026-03-16 14:29:13.613624	\N
\.


--
-- Name: room_participants_turn_order_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.room_participants_turn_order_seq', 1, false);


--
-- Name: room_participants room_participants_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.room_participants
    ADD CONSTRAINT room_participants_pkey PRIMARY KEY (room_id, user_id);


--
-- Name: rooms rooms_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rooms
    ADD CONSTRAINT rooms_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: vocabularies vocabularies_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.vocabularies
    ADD CONSTRAINT vocabularies_pkey PRIMARY KEY (language);


--
-- Name: room_participants room_participants_room_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.room_participants
    ADD CONSTRAINT room_participants_room_id_fkey FOREIGN KEY (room_id) REFERENCES public.rooms(id);


--
-- Name: room_participants room_participants_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.room_participants
    ADD CONSTRAINT room_participants_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: rooms rooms_admin_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rooms
    ADD CONSTRAINT rooms_admin_fkey FOREIGN KEY (admin) REFERENCES public.users(id);


--
-- Name: rooms rooms_current_player_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rooms
    ADD CONSTRAINT rooms_current_player_id_fkey FOREIGN KEY (current_player_id) REFERENCES public.users(id);


--
-- Name: rooms rooms_language_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rooms
    ADD CONSTRAINT rooms_language_fkey FOREIGN KEY (language) REFERENCES public.vocabularies(language);


--
-- Name: TABLE room_participants; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.room_participants TO alias;


--
-- Name: SEQUENCE room_participants_turn_order_seq; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON SEQUENCE public.room_participants_turn_order_seq TO alias;


--
-- Name: TABLE rooms; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.rooms TO alias;


--
-- Name: TABLE users; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.users TO alias;


--
-- Name: TABLE vocabularies; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.vocabularies TO alias;


--
-- Name: DEFAULT PRIVILEGES FOR SEQUENCES; Type: DEFAULT ACL; Schema: public; Owner: postgres
--

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON SEQUENCES TO alias;


--
-- Name: DEFAULT PRIVILEGES FOR FUNCTIONS; Type: DEFAULT ACL; Schema: public; Owner: postgres
--

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON FUNCTIONS TO alias;


--
-- Name: DEFAULT PRIVILEGES FOR TABLES; Type: DEFAULT ACL; Schema: public; Owner: postgres
--

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON TABLES TO alias;


--
-- PostgreSQL database dump complete
--

\unrestrict hjMldxUvJNq16ROR1Ywn56I5d8qbxk2RKLJkXnC3C5ioF0ExZeqFkUl2JPZPXUk

