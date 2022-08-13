<template>
  <accordion :multiple="true">
    <accordion-tab>
      <template #header>
        <div style="white-space: nowrap">
          The distinction between the <i>responsible</i> and <i>supporting</i> roles
        </div>
      </template>
      <p>For most purposes, the <i>responsible</i> and <i>supporting</i> roles are the same.
      Each of them equally contributes to the solution, and each of them is expected to bring
      a point of view, defend it, and listen to others as part of refining the solution.</p>
      <p>But every team needs a leader. That leader, in addition to all of the above, needs to
      make sure that people are participating in the decision making process in a respectful
      manner. They also need to deal with unresolved differences in opinion. That leadership
      role can't be shared, for risk of deadlock or thrashing on decisions.</p>
      <p>Separating the <i>responsible</i> from the <i>supporting</i> roles takes those problem on.
        It ensures that the <i>accountable</i> individual can always know from whom to get status
        reports, and eliminates finger pointing in terms of knowing who must ensure progress is made.</p>
    </accordion-tab>
    <accordion-tab>
      <template #header>
        <div style="white-space: nowrap">
          The <i>accepting</i> role
        </div>
      </template>
      <p>The <i>accepting</i> role is unique to ARSCIV-X. That uniqueness warrants a bit of
      discussion and justification.</p>
      <p>
        The <i>accepting</i> role was created after an analysis of why architectural decisions often
        languish or outright fail. In the end, architects rarely implement decisions all on their own.
        The intent of architecture decisions is getting to a common approach; the more individuals are
        involved in the implementation of that approach the more need there is for ensuring that there
        is not just awareness but also understanding of what the decisions say, and why they say it.
      </p>
      <p>
        Equally, there needs to be a guard against architectural decisions becoming too "ivory tower" -
        that is to say architecturally perfect but unimplementable given organizational or other
        technology constraints.
      </p>
      <p>The <i>accepting</i> role provides structural support for that check and balance. The
      <i>accepting</i> individual(s) do not have to be part of the day-to-day crafting of a decision, but
      do need to sign off the decision as being implementable. More importantly, as part of the
      Architecture Decision Process they not only sign off that it is implementable but <i>commit to
        actually implement it</i>.</p>
    </accordion-tab>
    <accordion-tab>
      <template #header>
        <div style="white-space: nowrap">
          Why the <i>accepting</i> role stands alone
        </div>
      </template>
      <p>One might wonder why the <i>accepting</i> role is separated from the other roles
        in the acronym. That is an explicit choice that emphasizes that there is often
        an organizational handoff as the output of the working group is finalized. For example,
        the architecture team will read out the results and handoff to the engineering team.</p>
      <p>This handoff is highlighted by explicitly separating out the role from the others.</p>
    </accordion-tab>
  </accordion>
</template>

<script>

export default {
  name: "ARSCIV_X_Definition"
}
</script>

<style scoped>

</style>