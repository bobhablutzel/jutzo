<template>
  <accordion :multiple="true">
    <accordion-tab>
      <template #header>
        <div style="white-space: nowrap">
          A two-pizza team
        </div>
      </template>
      <p>By design two-pizza teams are small, and so there will be a number of individuals that
      hold multiple roles. Consider a team consisting of a product owner, engineering lead, and
      engineers. In this case, you might have a structure something like this:
      </p>
      <ul>
        <li><i>Product Owner</i>: Accountable, Accepting</li>
        <li><i>Engineering Lead</i>: Responsible, Validating</li>
        <li><i>Engineer 1</i>: Supporting</li>
        <li><i>Engineer 2</i>: Supporting</li>
        <li><i>Engineer 3</i>: Informed</li>
        <li><i>Engineering Lead, sibling team A</i>: Consulted</li>
        <li><i>Engineering Lead, sibling team B</i>: Informed</li>
      </ul>
    </accordion-tab>
    <accordion-tab>
      <template #header>
        <div style="white-space: nowrap">
          A domain-level working group
        </div>
      </template>
      <p>In this example, a domain team is looking to create a new architectural standard for their
      product. This product will be used by other teams within the organization, though those clients
      will largely be using the product through a well defined contract. At the same time, the product
      needs to align with the enterprise standards around security, reliability, and other concerns. The
        domain team has a dedicated architect, but the architect is more junior and has not run a working
        group before. In light of this, the domain team has asked for help from the central architecture team.
      </p>
      <ul>
        <li><i>Product Owner</i>: Accountable, Informed</li>
        <li><i>Central Architect</i>: Responsible</li>
        <li><i>Domain Architect</i>: Supporting, Validating</li>
        <li><i>Domain Engineering Lead</i>: Consulted, Accepting</li>
        <li><i>Enterprise Security Architect</i>: Architect</li>
        <li><i>Client Team Architect</i>: Supporting</li>
        <li><i>Sr. Domain Engineer</i>: Consulted</li>
        <li><i>Domain Engineer</i>: Informed</li>
      </ul>
    </accordion-tab>
\  </accordion>
</template>

<script>

export default {
  name: "ARSCIV_X_Definition"
}
</script>

<style scoped>

</style>