package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/secretnamebasis/secret-site/app"
	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/models"
)

// SERVER

// pre-conditions
const ( // Endpoint configuration
	endpoint      = "http://127.0.0.1:3000/api"
	routeApiUsers = "/users"
	routeApiItems = "/items"
	user          = "secret"
	dbPath        = "./database/"
)

// Define testCase type
type testCase struct {
	name string
	fn   func(*testing.T)
}

// Define response type
type response struct {
	Data   interface{} `json:"data"`
	Status string      `json:"status"`
}

// test-conditions
var (
	// Configure server settings
	c = config.Server{
		Port: 3000,
		Env:  "testing",
	}

	delay = 1 * time.Nanosecond

	// Item Test Data
	//
	// we don't store empty content
	// Fail cases
	failItemCreateData = models.Item{
		Title:   "title",
		Content: "",
	}
	// we don't store empty titles
	failItemUpdateData = models.Item{
		Title:   "",
		Content: "content",
	}
	// Success cases
	successItemCreateData = models.Item{
		Title: "title",
		// love you Joyce
		Content: " no thats no way for him has he no manners nor no refinement nor no nothing in his nature slapping us behind like that on my bottom because I didnt call him Hugh the ignoramus that doesnt know poetry from a cabbage thats what you get for not keeping them in their proper place pulling off his shoes and trousers there on the chair before me so barefaced without even asking permission and standing out that vulgar way in the half of a shirt they wear to be admired like a priest or a butcher or those old hypocrites in the time of Julius Caesar of course hes right enough in his way to pass the time as a joke sure you might as well be in bed with what with a lion God Im sure hed have something better to say for himself an old Lion would O well I suppose its because they were so plump and tempting in my short petticoat he couldnt resist they excite myself sometimes its well for men all the amount of pleasure they get off a womans body were so round and white for them always I wished I was one myself for a change just to try with that thing they have swelling up on you so hard and at the same time so soft when you touch it my uncle John has a thing long I heard those cornerboys saying passing the comer of Marrowbone lane my aunt Mary has a thing hairy because it was dark and they knew a girl was passing it didnt make me blush why should it either its only nature and he puts his thing long into my aunt Marys hairy etcetera and turns out to be you put the handle in a sweepingbrush men again all over they can pick and choose what they please a married woman or a fast widow or a girl for their different tastes like those houses round behind Irish street no but were to be always chained up theyre not going to be chaining me up no damn fear once I start I tell you for their stupid husbands jealousy why cant we all remain friends over it instead of quarrelling her husband found it out what they did together well naturally and if he did can he undo it hes coronado anyway whatever he does and then he going to the other mad extreme about the wife in Fair Tyrants of course the man never even casts a 2nd thought on the husband or wife either its the woman he wants and he gets her what else were we given all those desires for Id like to know I cant help it if Im young still can I its a wonder Im not an old shrivelled hag before my time living with him so cold never embracing me except sometimes when hes asleep the wrong end of me not knowing I suppose who he has any man thatd kiss a womans bottom Id throw my hat at him after that hed kiss anything unnatural where we havent 1 atom of any kind of expression in us all of us the same 2 lumps of lard before ever Id do that to a man pfooh the dirty brutes the mere thought is enough I kiss the feet of you senorita theres some sense in that didnt he kiss our halldoor yes he did what a madman nobody understands his cracked ideas but me still of course a woman wants to be embraced 20 times a day almost to make her look young no matter by who so long as to be in love or loved by somebody if the fellow you want isnt there sometimes by the Lord God I was thinking would I go around by the quays there some dark evening where nobodyd know me and pick up a sailor off the sea thatd be hot on for it and not care a pin whose I was only do it off up in a gate somewhere or one of those wildlooking gipsies in Rathfarnham had their camp pitched near the Bloomfield laundry to try and steal our things if they could I only sent mine there a few times for the name model laundry sending me back over and over some old ones odd stockings that blackguardlooking fellow with the fine eyes peeling a switch attack me in the dark and ride me up against the wall without a word or a murderer anybody what they do themselves the fine gentlemen in their silk hats that K C lives up somewhere this way coming out of Hardwicke lane the night he gave us the fish supper on account of winning over the boxing match of course it was for me he gave it I knew him by his gaiters and the walk and when I turned round a minute after just to see there was a woman after coming out of it too some filthy prostitute then he goes home to his wife after that only I suppose the half of those sailors are rotten again with disease O move over your big carcass out of that for the love of Mike listen to him the winds that waft my sighs to thee so well he may sleep and sigh the great Suggester Don Poldo de la Flora if he knew how he came out on the cards this morning hed have something to sigh for a dark man in some perplexity between 2 7s too in prison for Lord knows what he does that I dont know and Im to be slooching around down in the kitchen to get his lordship his breakfast while hes rolled up like a mummy will I indeed did you ever see me running Id just like to see myself at it show them attention and they treat you like dirt I dont care what anybody says itd be much better for the world to be governed by the women in it you wouldnt see women going and killing one another and slaughtering when do you ever see women rolling around drunk like they do or gambling every penny they have and losing it on horses yes because a woman whatever she does she knows where to stop sure they wouldnt be in the world at all only for us they dont know what it is to be a woman and a mother how could they where would they all of them be if they hadnt all a mother to look after them what I never had thats why I suppose hes running wild now out at night away from his books and studies and not living at home on account of the usual rowy house I suppose well its a poor case that those that have a fine son like that theyre not satisfied and I none was he not able to make one it wasnt my fault we came together when I was watching the two dogs up in her behind in the middle of the naked street that disheartened me altogether I suppose I oughtnt to have buried him in that little woolly jacket I knitted crying as I was but give it to some poor child but I knew well Id never have another our 1st death too it was we were never the same since O Im not going to think myself into the glooms about that any more I wonder why he wouldnt stay the night I felt all the time it was somebody strange he brought in instead of roving around the city meeting God knows who nightwalkers and pickpockets his poor mother wouldnt like that if she was alive ruining himself for life perhaps still its a lovely hour so silent I used to love coming home after dances the air of the night they have friends they can talk to weve none either he wants what he wont get or its some woman ready to stick her knife in you I hate that in women no wonder they treat us the way they do we are a dreadful lot of bitches I suppose its all the troubles we have makes us so snappy Im not like that he could easy have slept in there on the sofa in the other room I suppose he was as shy as a boy he being so young hardly 20 of me in the next room hed have heard me on the chamber arrah what harm Dedalus I wonder its like those names in Gibraltar Delapaz Delagracia they had the devils queer names there father Vilaplana of Santa Maria that gave me the rosary Rosales y OReilly in the Calle las Siete Revueltas and Pisimbo and Mrs Opisso in Governor street O what a name Id go and drown myself in the first river if I had a name like her O my and all the bits of streets Paradise ramp and Bedlam ramp and Rodgers ramp and Crutchetts ramp and the devils gap steps well small blame to me if I am a harumscarum I know I am a bit I declare to God I dont feel a day older than then I wonder could I get my tongue round any of the Spanish como esta usted muy bien gracias y usted see I havent forgotten it all I thought I had only for the grammar a noun is the name of any person place or thing pity I never tried to read that novel cantankerous Mrs Rubio lent me by Valera with the questions in it all upside down the two ways I always knew wed go away in the end I can tell him the Spanish and he tell me the Italian then hell see Im not so ignorant what a pity he didnt stay Im sure the poor fellow was dead tired and wanted a good sleep badly I could have brought him in his breakfast in bed with a bit of toast so long as I didnt do it on the knife for bad luck or if the woman was going her rounds with the watercress and something nice and tasty there are a few olives in the kitchen he might like I never could bear the look of them in Abrines I could do the criada the room looks all right since I changed it the other way you see something was telling me all the time Id have to introduce myself not knowing me from Adam very funny wouldnt it Im his wife or pretend we were in Spain with him half awake without a Gods notion where he is dos huevos estrellados senor Lord the cracked things come into my head sometimes itd be great fun supposing he stayed with us why not theres the room upstairs empty and Millys bed in the back room he could do his writing and studies at the table in there for all the scribbling he does at it and if he wants to read in bed in the morning like me as hes making the breakfast for 1 he can make it for 2 Im sure Im not going to take in lodgers off the street for him if he takes a gesabo of a house like this Id love to have a long talk with an intelligent welleducated person Id have to get a nice pair of red slippers like those Turks with the fez used to sell or yellow and a nice semitransparent morning gown that I badly want or a peachblossom dressing jacket like the one long ago in Walpoles only 8/6 or 18/6 Ill just give him one more chance Ill get up early in the morning Im sick of Cohens old bed in any case I might go over to the markets to see all the vegetables and cabbages and tomatoes and carrots and all kinds of splendid fruits all coming in lovely and fresh who knows whod be the 1st man Id meet theyre out looking for it in the morning Mamy Dillon used to say they are and the night too that was her massgoing Id love a big juicy pear now to melt in your mouth like when I used to be in the longing way then Ill throw him up his eggs and tea in the moustachecup she gave him to make his mouth bigger I suppose hed like my nice cream too I know what Ill do Ill go about rather gay not too much singing a bit now and then mi fa pieta Masetto then Ill start dressing myself to go out presto non son piu forte Ill put on my best shift and drawers let him have a good eyeful out of that to make his micky stand for him Ill let him know if thats what he wanted that his wife is fucked yes and damn well fucked too up to my neck nearly not by him 5 or 6 times handrunning theres the mark of his spunk on the clean sheet I wouldnt bother to even iron it out that ought to satisfy him if you dont believe me feel my belly unless I made him stand there and put him into me Ive a mind to tell him every scrap and make him do it out in front of me serve him right its all his own fault if I am an adulteress as the thing in the gallery said O much about it if thats all the harm ever we did in this vale of tears God knows its not much doesnt everybody only they hide it I suppose thats what a woman is supposed to be there for or He wouldnt have made us the way He did so attractive to men then if he wants to kiss my bottom Ill drag open my drawers and bulge it right out in his face as large as life he can stick his tongue 7 miles up my hole as hes there my brown part then Ill tell him I want £ 1 or perhaps 30/- Ill tell him I want to buy underclothes then if he gives me that well he wont be too bad I dont want to soak it all out of him like other women do I could often have written out a fine cheque for myself and write his name on it for a couple of pounds a few times he forgot to lock it up besides he wont spend it Ill let him do it off on me behind provided he doesnt smear all my good drawers O I suppose that cant be helped Ill do the indifferent 1 or 2 questions Ill know by the answers when hes like that he cant keep a thing back I know every turn in him Ill tighten my bottom well and let out a few smutty words smellrump or lick my shit or the first mad thing comes into my head then Ill suggest about yes O wait now sonny my turn is coming Ill be quite gay and friendly over it O but I was forgetting this bloody pest of a thing pfooh you wouldnt know which to laugh or cry were such a mixture of plum and apple no Ill have to wear the old things so much the better itll be more pointed hell never know whether he did it or not there thats good enough for you any old thing at all then Ill wipe him off me just like a business his omission then Ill go out Ill have him eying up at the ceiling where is she gone now make him want me thats the only way a quarter after what an unearthly hour I suppose theyre just getting up in China now combing out their pigtails for the day well soon have the nuns ringing the angelus theyve nobody coming in to spoil their sleep except an odd priest or two for his night office or the alarmclock next door at cockshout clattering the brains out of itself let me see if I can doze off 1 2 3 4 5 what kind of flowers are those they invented like the stars the wallpaper in Lombard street was much nicer the apron he gave me was like that something only I only wore it twice better lower this lamp and try again so as I can get up early Ill go to Lambes there beside Findlaters and get them to send us some flowers to put about the place in case he brings him home tomorrow today I mean no no Fridays an unlucky day first I want to do the place up someway the dust grows in it I think while Im asleep then we can have music and cigarettes I can accompany him first I must clean the keys of the piano with milk whatll I wear shall I wear a white rose or those fairy cakes in Liptons I love the smell of a rich big shop at 7 1/2d a lb or the other ones with the cherries in them and the pinky sugar 11d a couple of lbs of those a nice plant for the middle of the table Id get that cheaper in wait wheres this I saw them not long ago I love flowers Id love to have the whole place swimming in roses God of heaven theres nothing like nature the wild mountains then the sea and the waves rushing then the beautiful country with the fields of oats and wheat and all kinds of things and all the fine cattle going about that would do your heart good to see rivers and lakes and flowers all sorts of shapes and smells and colours springing up even out of the ditches primroses and violets nature it is as for them saying theres no God I wouldnt give a snap of my two fingers for all their learning why dont they go and create something I often asked him atheists or whatever they call themselves go and wash the cobbles off themselves first then they go howling for the priest and they dying and why why because theyre afraid of hell on account of their bad conscience ah yes I know them well who was the first person in the universe before there was anybody that made it all who ah that they dont know neither do I so there you are they might as well try to stop the sun from rising tomorrow the sun shines for you he said the day we were lying among the rhododendrons on Howth head in the grey tweed suit and his straw hat the day I got him to propose to me yes first I gave him the bit of seedcake out of my mouth and it was leapyear like now yes 16 years ago my God after that long kiss I near lost my breath yes he said I was a flower of the mountain yes so we are flowers all a womans body yes that was one true thing he said in his life and the sun shines for you today yes that was why I liked him because I saw he understood or felt what a woman is and I knew I could always get round him and I gave him all the pleasure I could leading him on till he asked me to say yes and I wouldnt answer first only looked out over the sea and the sky I was thinking of so many things he didnt know of Mulvey and Mr Stanhope and Hester and father and old captain Groves and the sailors playing all birds fly and I say stoop and washing up dishes they called it on the pier and the sentry in front of the governors house with the thing round his white helmet poor devil half roasted and the Spanish girls laughing in their shawls and their tall combs and the auctions in the morning the Greeks and the jews and the Arabs and the devil knows who else from all the ends of Europe and Duke street and the fowl market all clucking outside Larby Sharons and the poor donkeys slipping half asleep and the vague fellows in the cloaks asleep in the shade on the steps and the big wheels of the carts of the bulls and the old castle thousands of years old yes and those handsome Moors all in white and turbans like kings asking you to sit down in their little bit of a shop and Ronda with the old windows of the posadas 2 glancing eyes a lattice hid for her lover to kiss the iron and the wineshops half open at night and the castanets and the night we missed the boat at Algeciras the watchman going about serene with his lamp and O that awful deepdown torrent O and the sea the sea crimson sometimes like fire and the glorious sunsets and the figtrees in the Alameda gardens yes and all the queer little streets and the pink and blue and yellow houses and the rosegardens and the jessamine and geraniums and cactuses and Gibraltar as a girl where I was a Flower of the mountain yes when I put the rose in my hair like the Andalusian girls used or shall I wear a red yes and how he kissed me under the Moorish wall and I thought well as well him as another and then I asked him with my eyes to ask again yes and then he asked me would I yes to say yes my mountain flower and first I put my arms around him yes and drew him down to me so he could feel my breasts all perfume yes and his heart was going like mad and yes I said yes I will Yes. ",
	}
	successItemUpdateData = models.Item{
		Title:   "squirrel",
		Content: "plain text",
	}

	// Fail cases
	// resopnse, err := dero.GetEncryptedBalance(address)
	// response.Result.Status != "OK"
	failCreateAddress  = "dero1qynmz4tgkmtmmspqmywvjjmtl0x8vn5ahz4xwaldw0hu6r5500hryqgptvnj"
	failUpdateAddress  = "dero1qyvqpdftj8r6005xs20rnflakmwa5pdxg9vcjzdcuywq2t8skqhvwqglt6x0"
	failUserCreateData = models.User{
		User:   user,
		Wallet: failCreateAddress,
	}
	failUserUpdateData = models.User{
		User:   user,
		Wallet: failUpdateAddress,
	}
	// Success cases
	successCreateAddress = "dero1qynmz4tgkmtmmspqmywvjjmtl0x8vn5ahz4xwaldw0hu6r5500hryqgptvnj8"
	successUpdateAddress = "dero1qyvqpdftj8r6005xs20rnflakmwa5pdxg9vcjzdcuywq2t8skqhvwqglt6x0g"

	successUserCreateData = models.User{
		User:   user,
		Wallet: successCreateAddress,
	}
	successUserUpdateData = models.User{
		User:   user,
		Wallet: successUpdateAddress,
	}

	// Test cases
	testCases     = append(itemTestCases, userTestCases...)
	itemTestCases = []testCase{
		// Item test cases
		{
			"CheckItems",
			checkItemsTest,
		},
		{
			"Create when Item is invalid",
			createItemFailTest,
		},
		{
			"Retrieve when Item 1 creation fails",
			retrieveItemFailTest,
		},
		{
			"Create when Item is valid",
			createItemSuccessTest,
		},
		{
			"Retrieve when Item 1 is successfully created",
			retrieveItemSuccessTest,
		},
		{
			"CheckItems",
			checkItemsTest,
		},
		{
			"Update when Item is not valid",
			updateItemFailTest,
		},
		{
			"Retrieve when Item 1 update fails",
			retrieveItemSuccessTest,
		},
		{
			"Update when Item is valid",
			updateItemSuccessTest,
		},
		{
			"Retrieve when Item 1 is successfully updated",
			retrieveItemSuccessTest,
		},
		{
			"Delete when Item 1 is present",
			deleteItemSuccessTest,
		},
		{
			"Delete when Item is not present",
			deleteItemFailTest,
		},
		{
			"Retrieve when Item 1 is deleted",
			retrieveItemFailTest,
		},
	}

	userTestCases = []testCase{
		// User test cases
		{
			"CheckUsers",
			checkUsersTest,
		},
		{
			"Create when user is invalid",
			createUserFailTest,
		},
		{
			"Retrieve when user 1 creation fails",
			retrieveUserFailTest,
		},
		{
			"Create when user is valid",
			createUserSuccessTest,
		},
		{
			"Retrieve when user 1 is successfully created",
			retrieveUserSuccessTest,
		},
		{
			"Update when user is invalid",
			updateFailTest,
		},
		{
			"Retrieve when user 1 update fails",
			retrieveUserSuccessTest,
		},
		{
			"Update when user is valid",
			updateUserSuccessTest,
		},
		{
			"Retrieve when user 1 is successfully updated",
			retrieveUserSuccessTest,
		},
		{
			"Delete when user1 is present",
			deleteUserSuccessTest,
		},
		{
			"Delete when user is not present",
			deleteUserFailTest,
		},
		{
			"Retrieve when user 1 is deleted",
			retrieveUserFailTest,
		},
	}
)

func TestApi(t *testing.T) {

	// Start the server and handle shutdown
	a := startServer()

	// Run tests
	runTests(t)

	// Stop the server after tests are done
	stopServer(t, a)

	// Delete the database
	deleteDB()
}

// test-server
func startServer() *app.App { // start the server
	// Delete the database before starting the server
	deleteDB()

	a := app.MakeApp(c)
	go func() {
		if err := a.StartApp(c); err != nil {
			log.Fatalf("Error starting server: %s\n", err)
		}
	}()
	return a
}

func stopServer(t *testing.T, a *app.App) { // stop the server
	// Stop the server after tests are done
	if err := a.StopApp(); err != nil {
		t.Errorf("Error stopping server: %s\n", err)
	} else {
		log.Println("Server stopped gracefully")
	}
}

func deleteDB() {
	err := os.RemoveAll(dbPath)
	if err != nil {
		log.Fatalf("Error deleting database: %s\n", err)
	}
	log.Println("Database deleted successfully")
}

// TEST

// Define a test execution
func executeTest(t *testing.T, actionFunc func() (string, error), expectedStatus string) {

	// Execute the action function
	responseBody, err := actionFunc()
	if err != nil {
		t.Fatalf("Error executing action: %v", err)
	}

	// Unmarshal the response body into the response struct
	var resp response
	if err := json.Unmarshal([]byte(responseBody), &resp); err != nil {
		t.Fatalf("Error parsing response: %v", err)
	}

	// Compare the expected status with the actual status
	if expectedStatus != resp.Status {
		t.Errorf("Expected status: %s, Actual status: %s", expectedStatus, resp.Status)
	}
	// log.Printf("%s", resp)
	// Sleep for 1 nanosecond
	time.Sleep(delay)
}

// performAction performs an HTTP request with the provided method, URL, and data.
func performAction(method, url string, data interface{}) (string, error) {

	// Marshal data into JSON payload
	payload, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON data: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return "", fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	return string(body), nil
}

// Run tests
func runTests(t *testing.T) {
	log.Printf("Environment: %s\n", c.Env)
	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.name, tc.fn)
	}
}

// ITEM
// Functions to perform Item API CRUD actions
func checkItems() (string, error) {
	return performAction(
		"GET",
		endpoint+routeApiItems,
		nil,
	)
}
func checkItemsTest(t *testing.T) {
	executeTest(t, checkItems, "success")
}

// CREATE

// CREATE FAIL
func createItemFail() (string, error) {
	return performAction("POST", endpoint+routeApiItems, failItemCreateData)
}
func createItemFailTest(t *testing.T) {
	executeTest(t, createItemFail, "error")
}

// CREATE SUCCESS
func createItemSuccess() (string, error) {
	return performAction("POST", endpoint+routeApiItems, successItemCreateData)
}
func createItemSuccessTest(t *testing.T) {
	executeTest(t, createItemSuccess, "success")
}

// RETREIVE
func retrieveItem() (string, error) {
	return performAction("GET", fmt.Sprintf("%s/1", endpoint+routeApiItems), nil)
}

// RETREIVE SUCCESS
func retrieveItemSuccessTest(t *testing.T) {
	executeTest(t, retrieveItem, "success")
}

// RETREIVE FAIL
func retrieveItemFailTest(t *testing.T) {
	executeTest(t, retrieveItem, "error")
}

// UPDATE

// UPDATE FAIL
func updateItemFail() (string, error) {
	return performAction("PUT", fmt.Sprintf("%s/1", endpoint+routeApiItems), failItemUpdateData)
}
func updateItemFailTest(t *testing.T) {
	executeTest(t, updateItemFail, "error")
}

// UPDATE SUCCESS
func updateItemSuccess() (string, error) {
	return performAction("PUT", fmt.Sprintf("%s/1", endpoint+routeApiItems), successItemUpdateData)
}
func updateItemSuccessTest(t *testing.T) {
	executeTest(t, updateItemSuccess, "success")
}

// DELETE
func deleteItem() (string, error) {
	return performAction("DELETE", fmt.Sprintf("%s/1", endpoint+routeApiItems), nil)
}

// DELETE FAIL
func deleteItemFailTest(t *testing.T) {
	executeTest(t, deleteItem, "error")
}

// DELETE SUCCESS
func deleteItemSuccessTest(t *testing.T) {
	executeTest(t, deleteItem, "success")
}

// USER

// Check all users
func checkUsers() (string, error) {
	return performAction("GET", endpoint+routeApiUsers, nil)
}
func checkUsersTest(t *testing.T) {
	executeTest(t, checkUsers, "success")
}

// CREATE

// CREATE FAIL
func createUserFail() (string, error) {
	return performAction("POST", endpoint+routeApiUsers, failUserCreateData)
}
func createUserFailTest(t *testing.T) {
	executeTest(t, createUserFail, "error")
}

// CREATE SUCCESS
func createUserSuccess() (string, error) {
	return performAction("POST", endpoint+routeApiUsers, successUserCreateData)
}
func createUserSuccessTest(t *testing.T) {
	executeTest(t, createUserSuccess, "success")
}

// RETREIVE
func retrieveUser() (string, error) {
	return performAction("GET", fmt.Sprintf("%s/1", endpoint+routeApiUsers), nil)
}

// RETREIVE SUCCESS
func retrieveUserSuccessTest(t *testing.T) {
	executeTest(t, retrieveUser, "success")
}

// RETREIVE FAIL
func retrieveUserFailTest(t *testing.T) {
	executeTest(t, retrieveUser, "error")
}

// UPDATE

// UPDATE FAIL
func updateUserFail() (string, error) {
	return performAction("PUT", fmt.Sprintf("%s/1", endpoint+routeApiUsers), failUserUpdateData)
}
func updateFailTest(t *testing.T) {
	executeTest(t, updateUserFail, "error")
}

// UPDATE SUCCESS
func updateUserSuccess() (string, error) {
	return performAction("PUT", fmt.Sprintf("%s/1", endpoint+routeApiUsers), successUserUpdateData)
}
func updateUserSuccessTest(t *testing.T) {
	executeTest(t, updateUserSuccess, "success")
}

// DELETE
func deleteUser() (string, error) {
	return performAction("DELETE", fmt.Sprintf("%s/1", endpoint+routeApiUsers), nil)
}

// DELETE FAIL
func deleteUserFailTest(t *testing.T) {
	executeTest(t, deleteUser, "error")
}

// DELETE SUCCESS
func deleteUserSuccessTest(t *testing.T) {
	executeTest(t, deleteUser, "success")
}
