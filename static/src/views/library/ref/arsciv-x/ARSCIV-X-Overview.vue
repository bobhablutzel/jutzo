<template>
  <div class="text-container">
    <div class="text-container">
      <p>Organizations often adopt
        <a href="https://en.wikipedia.org/wiki/Responsibility_assignment_matrix" target="_blank">
          responsibility assignment matrices</a>, often known as RACI models, in order to articulate the
        roles and responsibilities of various stakeholders in the organization. While wholeheartedly supporting
        these models, it's also important to recognize that the simpler models don't often provide
        enough clarity for large-scale architecture practices. In particular, they</p>
      <ul>
        <li>Don't account for the necessary handoff between architecture and engineering; and</li>
        <li>Risk creating confusion on who is leading the effort</li>
      </ul>
      <p>One can adopt an <i>ARSCIV-X</i> model to address these short comings. In the ARSCIV-X
        approach, there are unique roles for the leader of an initiative vs the participants, and
        for the individuals required to sign off on the work before it is considered complete.
      </p>
      <p>
        While general purpose, ARSCIV-X is specifically designed for modern engineering teams. Many organizations
        adopt a 3-in-a-box model for technology teams, where the product owner, architect, and engineering
        lead represent their individual concerns as a team. In this model, each of the individuals
        has responsibility in their area but also, because of the inter-dependencies, have to agree on
        the overall outputs. So an architect might have the lead for making an architectural decision but needs
        to ensure that it is still aligned with the product goals and actually implementable.
      </p>
      <p>
        In small teams this occurs relatively naturally; in the smallest teams there might not even be a well
        defined distinction between some of the roles. (Often you'll see architects and engineers combined.) For
        larger teams, or ones that require more visibility into the decision making process, the rigor of a
        responsibility assignment matrix becomes more valuable.
      </p>
      <p>
        ARSCIV-X is unique in that it acknowledges the need for a distinction between the accountable / responsible
        parties and the ones that will have to live with the decisions being made. It formalizes the need for the
        engineering and product leads to agree with the outputs, and gives them a specific role in the overall
        responsibility matrix.
      </p>
      <p>
        ARSCIV-X is also unique in that it explicitly requires a single leader for an effort. This prevents the
        confusion of simpler RACI models where multiple people might be marked as "responsible" for getting something
        done. That's usually functionally true, but creates uncertainty about who is driving the initiative and
        ensuring that work is done in a timely manner. Worse, it prevents there being a single person who can be
        evaluated on whether or not work is being done. The separation of the "responsible" role - which <b>MUST</b> be a
        single person - from the "supporting" role - which equally contributes but does not drive - is another
        differentiating characteristic of the ARSCIV-X model.
      </p>
    </div>

    <p>The ARSCIV-X model is an important part of the Architecture Decision Process, which formalizes the hand-off
      between the responsible and accepting roles. </p>

  </div>
</template>

<script>

export default {
  name: "ARSCIV_X_Overview"
}
</script>

<style scoped>

</style>