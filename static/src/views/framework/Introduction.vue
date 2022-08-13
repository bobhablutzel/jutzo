<template>
  <panel header="Scalable Enterprise Architecture Framework">
    <div class="text-container">
      <img class="frameworkAbstract" :src="'/images/EA Framework Abstract.svg'" alt="Abstract EA Framework"/>
      <p>
        Businesses globally spend billions of dollars annually on information technology,
        often with disappointing results. Not disappointing results as in “we missed our deadline”,
        but disappointing results as in “we are not getting a significant uplift in business results”.
        Instead, companies find themselves mired in tech debt, beholden to ever-expanding technology
        organizations, and constrained by an inability to bring new products to market.
      </p>
      <p>
        For sure, it isn’t always this way. Some companies excel at using technology to strategic advantage.
        Leading companies such as Capital One, PayPal, Venmo, Google, Amazon, and others are
        successfully building business systems that are cost effective, business focused, and game changing.
        These companies enjoy the benefit of not only reducing their operating costs but being able to rapidly
        innovate and iterate on new products – lapping their competitors in the marketplace.
        They enjoy greater profit margins and revenues, but more importantly they delight their
        customers and build enduring brands.
      </p>
      <p>
        I’ve spent much of my adult life working at and studying these companies. Specifically, I’ve been
        fascinated by the companies that use technology as a foundational part of the business and make the
        translation from business needs to technology implementations seamless. Of course, I’ve equally studied
        companies that fail at this – sometimes spectacularly – and learned as much from the failures of these
        organizations as the success of others.
      </p>
      <p>
        Enterprise Architecture (EA) is often touted as the solution to the problem – a structured approach
        to tying business strategy to technology. Realistically though, many EA programs fail, or because academic
        exercises that are disconnected from either the business or the technologists – or both.
      </p>
      <p>
        This problem has particularly fascinated me. If the whole point of Enterprise Architecture is to tie
        business to engineering, why does it so frequently miss the mark?
      </p>
      <p>
        Let’s be realistic about two things.
      </p>
      <ol>
        <li>The value of Enterprise Architecture increases with the complexity of the organization. The whole premise of
          EA is that you need some centralized body that is making sure all the pieces are coming together in a
          manageable
          and efficient way so that you can avoid the duplication and unnecessary work that drives up the Run the Engine
          expenses. Smaller companies simply don’t need as robust an EA function; it is easier for the individual teams
          to
          track what each other are doing and avoid that duplication organically.
        </li>
        <li>
          The effectiveness of Enterprise Architecture decreases with the complexity of the organization. Centralizing
          the Enterprise Architecture practice creates a bottleneck in the organization, impeding delivery times and
          stifling innovation. The only way for a centralized team to keep pace with large delivery organizations is to
          scale to impractical levels. Either the company cannot sustain the investment, or the sheer size of the EA
          team
          creates divergence within the team itself. Either way, the practical value of Enterprise Architecture is
          diminished simply by the inability to keep pace.
        </li>
      </ol>
      <p>
        This creates what I call “the Pendulum” – a repeated pattern of centralization followed by disbanding followed
        by centralization – that may be familiar to you. Observing the Pendulum occur at multiple companies was the
        genesis of the SEA framework. It became almost a given that this cycle would occur at about 3-5 year intervals
        in most companies, and I wanted to know why.
      </p>
      <p>
        Thinking about this problem over the years led me to the realization that a new approach was in order – one that
        scaled with the complexity of the organization. My goals were simple:
      </p>
      <ul>
        <li>A framework for delivering on the linkage of business to technology in a practical way</li>
        <li>Ensuring that framework scaled for the most complex companies</li>
        <li>Assuring that framework was deployable without significant disruption to any existing EA initiatives</li>
      </ul>
      <p>
        The end result of that was the Scalable Enterprise Architecture (SEA) Framework, which this site describes in
        detail.
      </p>
    </div>
  </panel>
</template>


<script>

export default {
  name: "FrameworkIntroduction"
}
</script>

<style scoped>


.frameworkAbstract {
  max-width: 300px;
}


</style>