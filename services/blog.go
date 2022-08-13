package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
)

type BlogTextSection struct {
	Ordinal int    `json:"ordinal"`
	Text    string `json:"text"`
}

type BlogImgSection struct {
	Ordinal   int    `json:"ordinal"`
	Source    string `json:"src"`
	Alternate string `json:"alt"`
}

type BlogHeaderSection struct {
	Ordinal int    `json:"ordinal"`
	Level   int    `json:"level"`
	Text    string `json:"header"`
}

type BlogSummary struct {
	ID              uuid.UUID `json:"id"`
	PublicationDate time.Time `json:"publicationDate"`
	Title           string    `json:"title"`
	Teaser          string    `json:"teaser"`
}

type BlogEntry struct {
	BlogSummary
	Body []any `json:"body"`
}

type ErrorMessage struct {
	Error string `json:"error"`
}

var blog1Time, _ = time.Parse("2006-01-02", "2018-09-04")
var blog1 = BlogEntry{
	BlogSummary{
		uuid.MustParse("1f4bba93-d06d-4c74-b905-53f19fc5550d"),
		blog1Time,
		"Becoming an architect",
		"As software engineers progress in their careers, they often wonder if architecture is a possible career path. But what does becoming an architect even mean?",
	},
	[]any{
		BlogTextSection{0, "One of the most frequent questions I get from young engineers is “how do I become an architect?”. There are a number of answers to this, but one keeps coming back to me: <i>you have to fall out of love with technology</i>. Let me explain."},
		BlogTextSection{1, "At some point in your career you have to decide if you're on the engineering track, the architect track, or the management track. The manager track is enough of a different beast that I'm not going to talk about it here. Instead, let's focus on the difference between the engineer track and the architect track."},
		BlogTextSection{2, "The old platitude is that architects think broadly, while engineers think deeply. Many companies use the “T” model; engineers are the stem (focused deeply on a specific technology) while architects were the crossbar (focused across a range of technologies). While that's frequently the end result, it's the effect and not the cause."},
		BlogTextSection{3, "The cause is that architects focus on the risks to the system, while engineers focus on the implementations of the systems. This distinction is sometimes fuzzy, and in small systems a given individual might well swap back and forth between those roles. As the system gets larger, though, it makes sense to have individuals focused on these tasks full time."},
		BlogTextSection{4, "For the architect this means two things. First of all, you have to be unbiased in your selection of the right technology. The technology that was right for the last project might not be appropriate for the next one. Yes, once you get to Turing completeness you can pretty much to anything with any technology. But platforms are written with specific problem spaces in mind, and aligning the technology to the project can vastly simplify things."},
		BlogTextSection{5, "Second of all, you have to trust your engineers to implement the system. This can be really hard for engineers transitioning into architects; after all they have a history of being rewarded for being deep in the specifics of implementations. Letting go of that feels unnatural."},
		BlogTextSection{6, "When I mentor aspiring architects, I talk about the “Rule of 10”. Any given problem can be implement in 100 different ways. Of those, 90 suck for one reason or another - too slow, too expensive, too unreliable. Of the 10 that remain any one is acceptable. You might have picked design #3 while your engineering team picked design #7. The trick to being an architect is understanding that #7 is just as good (maybe even better by some measures) and letting go of the engineering micro-management. [Needless to say, if your engineering team suggests #87 you have to help them understand why it won't work.]"},
		BlogTextSection{7, "Both the unbiased view of technology and the Rule of 10 mean that you are letting go of your specific history and input in engineering. That's really hard to do if you love the technology. Loving the technology clouds your judgement, and drives you into wrong decisions and micro-management. Think of all the architects you hated to work for - in general it was because they ran roughshod over your technical abilities because they felt theirs were superior."},
		BlogTextSection{8, "[By the way, that in no way suggests that architects should become ivory tower, disconnected, or rusty. Keeping your technology skills sharp is a requirement for being an architect. But there is a difference between keeping up and keeping control, and the architect needs to let go if they are going to be successful.]"},
		BlogTextSection{9, "This is hard for some engineers; they really and truly love the hands-on experience of wrangling an elegant implementation to a given problem. That's fine. Architect isn't a promotion path, it's a career path that is separate but equal to engineer. When we wrote the upgraded career ladder for my current company we kept those two roles separate, and have individuals at all levels of the company on each of the career paths. Companies make the mistake that architect is a promotion path because in general you need fewer of them, but that's just not so."},
		BlogTextSection{10, "In the end the litmus test for whether to pursue the architect path is whether you love the technology, or whether you're willing to cede control of that while you focus on more abstract concerns (risks to the system). There's no right or wrong in this; you need to be true to yourself so that you don't hate going to work every day."},
		BlogTextSection{11, "Falling out of love with technology, while still keeping up to date with technology, is no easy task. The best architects I know are pragmatists who think of technology as tools, and can quickly assess overall strengths and weaknesses. Engineers tend to love exploring the nooks and crannies and learning how to make a given platform do things in a better, faster, or simply different way. Both are needed; both are important; but recognizing which one appeals to you more is critical to deciding what path you're on."},
	},
}

var blog2Time, _ = time.Parse("2006-01-02", "2018-07-31")
var blog2 = BlogEntry{
	BlogSummary{
		uuid.MustParse("aff73e3c-3142-4955-b76a-8b7a593f87bb"),
		blog2Time,
		"Defining Architecture",
		"After several years of being an architect, I finally have the “What is an architect” elevator pitch",
	},
	[]any{
		BlogTextSection{0, "I’ve been a practicing architect (in one way or another) for over 20 years. But “architecture” is notoriously hard to define – there are dozens of “architectures” and “architects” in any large system. Surely these must have a common thread – some way of talking about all these architectures in a common way – but we have to define what that is. "},
		BlogTextSection{1, "What I’ve came to realize – slowly over time – is this fundamental question of what we do is not commonly asked. Most architecture methodologies talk about the how – how to create models, how to document systems – and sometimes the what – whether that’s people, processes and things or entities and relationships – but not the why of architecture. Why is architecture important, why does it have a role to play, and why should it be a part of our project plans."},
		BlogTextSection{2, "What is architecture in the first place?"},
		BlogTextSection{3, "I started thinking a lot about this question a year or so ago. Embarrassingly, after 20 some years of being “an architect” I didn’t have a clean, concise definition of what architecture was. I didn’t have that “elevator pitch” I could use with an executive who asked what I did for the company. I certainly didn’t have an answer for my grandmother at Thanksgiving – at 102 years old, telling her I “made the world a better place through elegant hierarchies for maximum code reuse and extensibility” was about as meaningful as reciting “The Jabberwocky”, and a good deal less satisfying."},
		BlogTextSection{4, "I believe in architecture. I have been an architect most of my adult life, and I firmly believe is is critical to any organization or system of reasonable size. I like to think I can do architecture and be an architect – but I had to face that hard realization that no professional can be successful if it cannot sell itself to others. There had to be a way of explaining what I did, and why what I did was important, to people who hadn’t spend years of their lives studying UML, TOGAF, or any other approach to architecture."},
		BlogTextSection{5, "I also came to realize that I had to stop thinking like a technologist. The reality is that we are part of a business. I had to articulate what my profession did for the business; how it advanced our corporate mission in a way business people could relate to. The hard reality is that something that is poorly understood is seen as an overhead at best and unnecessary at worst and probably is. If architecture didn’t have a business function to fulfill and fulfill in a unique way that no other part of the business could, then yes, it probably is overhead and probably should be eliminated."},
		BlogTextSection{6, "The good news, though, is that architect does have that business function and more interestingly it can be summarized in a single phrase for all aspects of architecture."},
		BlogHeaderSection{7, 2, "The Pitch"},
		BlogTextSection{8, "I came to the conclusion that architecture controls against <i>quality risks to a system</i>. OK, maybe that’s not much better than “The Jabberwocky”, but it’s a starting point. It does make a couple things clear right off the bat. It makes clear that architecture doesn’t deal with functionality per se. As an architect, I’m not defining what the system does, I am defining how it does it. Functionality – what the customers get for using a system – is the purview of product management. They talk with customers, decide what customers need, and ruthlessly define the timelines against which to deliver those features to the customer – after all, customers are notoriously ADHD and will wander off into the wilderness of competitive products if we don’t constantly keep their attention. We, the architects and engineers, have to deliver those timelines – and make sure we don’t break anything in the process."},
		BlogTextSection{9, "It’s the not breaking things that has always fascinated me. It’s one thing to make a program work. It’s another thing completely to make it work relentlessly – every single time the customer makes a request, in a matter of milliseconds, millions of times a day, all while making sure no data is lost, corrupted, or accidentally disclosed. Building and reinforcing our customers’ trust."},
		BlogTextSection{10, "This crystallized as the essence of architecture for me. Product management defined what the system did. Engineering make sure it was implemented correctly. Architecture – my long undefined career – made sure the engineers knew what “correctly” meant, thinking of all the details that might go wrong and putting a plan in place to make sure that those conditions didn’t arise or were caught and corrected when they did, so that the customer could have a relentlessly enjoyable experience with the system."},
		BlogTextSection{11, "Of course, this definition has at least two problems. First, it could be too ivory tower – on the surface it’s somewhat hard to reconcile architectures like Lambda or Hadoop with “controls against quality risks”. Second, it uses a bunch of words that aren’t obviously defined: control, quality risk, and system. So in order to make this usable, let’s first get precise about these terms. Then we can get back to the question of whether this even makes sense in the first place. If it doesn’t – and, spoiler alert: I think it does – at least we’ll have the satisfaction of defining a couple terms. Defining precise terms is something I’ve always found architects to enjoy."},
		BlogHeaderSection{12, 3, "System"},
		BlogTextSection{13, "Let’s start with the easy one: system. A quick web search defines “system” as ”a set of connected things or parts forming a complex whole”. At least intuitively that makes sense. The more complicated our software and hardware environment gets, the more we need “architecture” and the more we worry about how the parts fit together."},
		BlogTextSection{14, "In our world those components can be a number of different things. They could be services as part of service oriented architecture. They could be networking, storage and compute equipment in our data centers, colos, or the cloud. The might be brands or subsidiary companies. The types of components determine the type of architecture (application architecture, network and data center architecture, and so forth). But in all cases we need a collection of connected, interoperating components in order to make architecture interesting; if your application can be written in a single Python module you might not have to think much about architecture at all."},
		BlogTextSection{15, "There also needs to be a common purpose to the system. Whether that’s routing data packets to processing units, or providing a merchant with a means to obtain checkout services, or just entertaining a kid by ringing bells with balls moving on wire tracks, the system needs to work toward a common goal. Just having a random collection of disconnected components doesn’t cut it – there needs to be some sense of order and purpose and integration between the components as well."},
		BlogTextSection{16, "So a system is a complex, interconnected set of components all working to achieve a common goal."},
		BlogHeaderSection{17, 3, "Quality Risk"},
		BlogTextSection{18, "“Quality Risk” is a bit trickier. What is ”quality” in the first place? “Quality” has a couple meanings in English; it can be “the degree of excellence of something” – how good or bad that something is when compared to other things of a similar kind. It can also be “a distinct attribute or characteristic” – for example, a leadership quality or that “indescribable quality that makes someone a star”. Of these definitions, we’re talking about the former: a measurement of how well something measures against a stated objective. It’s confusing because we sometimes talk about the “system qualities” meaning the second definition, as opposed to the “system quality” meaning the first. English is lovely sometimes."},
		BlogTextSection{19, "It’s important to notice that this almost always means a non-functional characteristic of a system. We don’t really speak about how well a system measures against the stated business objective; it either meets the objective, or it doesn’t. We can build an absolutely fabulous system for selling buggy whips, and probably still have a business failure even though the system quality was very high. This goes back to the discussion on product owners; they have to make sure that someone really wants the system in the first place. Our job as architects is to make sure it’s high quality if they do."},
		BlogTextSection{20, "We’ll talk more about this another time, but architecture quality can’t be subjective. “Quality” isn’t like “beauty”; it can’t be in the eye of the beholder. We need some way of specifying what our quality target is in clear, objective language. Something like “we will respond to the user within 300ms”. That’s pretty clear, and easy to measure. It also doesn’t say anything about how secure our response will be, or how much it will cost us to make that response. So we’ll have to have targets for those as well, and a dozen other things that make up the quality of our system."},
		BlogTextSection{21, "The risk, then is that we state we’ll make a certain measure and not actually make is. Maybe we respond in 301ms instead – not such a big deal. 3000ms, on the hand, could be and 30000ms definitely is. A quality risk, then is the risk that we might not make the target or, more precisely, that we’ll miss the target by a noticeable amount. That’s really what architecture focuses on; risks that systems won’t meet non-functional targets by a noticeable amount."},
		BlogHeaderSection{22, 3, "Control"},
		BlogTextSection{23, "That brings us to “control”. Since we’re taking about quality risks, it makes sense to talk about “control” in the language of risk management. A control in risk management lingo is “the policy or procedure by which potential risks are evaluated, reduced, or eliminated”. Controls are born out of thinking about what could go wrong, planning for what to do if they do go wrong, and planning on how to prevent them from going wrong in the first place."},
		BlogTextSection{24, "You’ll hear architects talking about the ”guardrails” of a system. That’s basically an informal way of saying ”control”. By putting in place a guardrail – which is to say by specifying a policy or procedures that we expect to be followed – we are creating an expectation of how the system will be designed in order to meet a set of quality goals. The controls don’t tell us exactly how to write a system, but they do keep us within a certain set of safe boundaries. A guardrail doesn’t tell us how to drive; it just tells us that straying too far from a certain path could be dangerous."},
		BlogTextSection{25, "Note also that the controls don’t say anything about what the system does. Guardrails keep you on the road, but don’t tell you where you’re going, or why you want to be there. Again the distinction between product management, and architecture."},
		BlogTextSection{26, "So now we can go back to the definition. Architecture is the result of someone thinking about what the desired characteristics of a system are, how to set a target for those characteristics, how to measure the actual system performance, and what policies and procedures to put in place to make achieving the target as likely as possible. We call those targets SLAs; we call those policies and procedures ”architecture”, and we call that someone an architect."},
	},
}

var blogs = []BlogEntry{blog1, blog2}

func newest(c *gin.Context) {
	log.Println("In newest function")
	var newestBlogs []BlogSummary
	for _, element := range blogs {
		newestBlogs = append(newestBlogs, element.BlogSummary)
	}
	c.JSON(http.StatusOK, newestBlogs)
}

func blogEntry(c *gin.Context) {
	id := c.Params.ByName("id")
	uniqueID, e := uuid.Parse(id)
	if e != nil {
		c.JSON(http.StatusBadRequest, ErrorMessage{"Invalid UUID"})
	} else {
		found := false
		for _, element := range blogs {
			if element.ID == uniqueID {
				c.JSON(http.StatusOK, element)
				found = true
			}
		}

		if !found {
			c.JSON(http.StatusNotFound, ErrorMessage{fmt.Sprintf("Blog id %s not found", id)})
		}
	}

}
