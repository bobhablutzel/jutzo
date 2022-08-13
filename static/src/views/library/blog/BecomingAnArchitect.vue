<template>
  <panel header="Becoming an Architect">
    <div class="text-container">
      <p>
        One of the most frequent questions I get from young engineers is "how do I become an architect?". There are a
        number of answers to this, but one keeps coming back to me: you have to fall out of love with technology. Let me
        explain.
      </p>
      <p>
        At some point in your career you have to decide if you're on the engineering track, the architect track, or the
        management track. The manager track is enough of a different beast that I'm not going to talk about it here.
        Instead, let's focus on the difference between the engineer track and the architect track.
      </p>
      <p>
        The old platitude is that architects think broadly, while engineers think deeply. At my previous company we used
        the "T" model; engineers were the stem (focused deeply on a specific technology) while architects were the
        crossbar (focused across a range of technologies). While that's frequently the end result, it's the effect and not
        the cause.
      </p>
      <p>
        The cause is that architects focus on the risks to the system, while engineers focus on the implementations of the
        systems. This distinction is sometimes fuzzy, and in small systems a given individual might well swap back and
        forth between those roles. As the system gets larger, though, it makes sense to have individuals focused on these
        tasks full time.
      </p>
      <p>
        For the architect this means two things. First of all, you have to be unbiased in your selection of the right
        technology. The technology that was right for the last project might not be appropriate for the next one. Yes,
        once you get to Turing completeness you can pretty much to anything with any technology. But platforms are written
        with specific problem spaces in mind, and aligning the technology to the project can vastly simplify things.
      </p>
      <p>
        Second of all, you have to trust your engineers to implement the system. This can be really hard for engineers
        transitioning into architects; after all they have a history of being rewarded for being deep in the specifics of
        implementations. Letting go of that feels unnatural.
      </p>
      <p>
        When I mentor aspiring architects, I talk about the "Rule of 10". Any given problem can be implement in 100
        different ways. Of those, 90 suck for one reason or another - too slow, too expensive, too unreliable. Of the 10
        that remain any one is acceptable. You might have picked design #3 while your engineering team picked design #7.
        The trick to being an architect is understanding that #7 is just as good (maybe even better by some measures) and
        letting go of the engineering micro-management. [Needless to say, if your engineering team suggests #87 you have
        to help them understand why it won't work.]
      </p>
      <p>
        Both the unbiased view of technology and the Rule of 10 mean that you are letting go of your specific history and
        input in engineering. That's really hard to do if you love the technology. Loving the technology clouds your
        judgement, and drives you into wrong decisions and micro-management. Think of all the architects you hated to work
        for - in general it was because they ran roughshod over your technical abilities because they felt theirs were
        superior.
      </p>
      <p>
        [By the way, that in no way suggests that architects should become ivory tower, disconnected, or rusty. Keeping
        your technology skills sharp is a requirement for being an architect. But there is a difference between keeping up
        and keeping control, and the architect needs to let go if they are going to be successful.]
      </p>
      <p>
        This is hard for some engineers; they really and truly love the hands-on experience of wrangling an elegant
        implementation to a given problem. That's fine. Architect isn't a promotion path, it's a career path that is
        separate but equal to engineer. When we wrote the upgraded career ladder for my current company we kept those two
        roles separate, and have individuals at all levels of the company on each of the career paths. Companies make the
        mistake that architect is a promotion path because in general you need fewer of them, but that's just not so.
      </p>
      <p>
        In the end the litmus test for whether to pursue the architect path is whether you love the technology, or whether
        you're willing to cede control of that while you focus on more abstract concerns (risks to the system). There's no
        right or wrong in this; you need to be true to yourself so that you don't hate going to work every day.
      </p>
      <p>
        Falling out of love with technology, while still keeping up to date with technology, is no easy task. The best
        architects I know are pragmatists who think of technology as tools, and can quickly assess overall strengths and
        weaknesses. Engineers tend to love exploring the nooks and crannies and learning how to make a given platform do
        things in a better, faster, or simply different way. Both are needed; both are important; but recognizing which
        one appeals to you more is critical to deciding what path you're on.
      </p>
    </div>
  </panel>
</template>


<script>

export default {
  name: "BecomingAnArchitect"
}
</script>

<style scoped>


</style>