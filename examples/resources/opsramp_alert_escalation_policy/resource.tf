resource "opsramp_alert_escalation_policy" "default_alert_escalation_policy" {
  name         = "Default escalation policy"
  tenant_scope = "CLIENT"
  precedence   = 1

  enabled_mode = "OBSERVED"

  escalation_type = "AUTOMATIC_UNTIL_ACKNOWLEDGED_CLOSED_SUPPRESSED_TICKETED"
  policy_type     = "ESCALATION_POLICY"

  escalations = [
    {
      wait_mins          = 0
      priority           = "Normal"
      repeat_frequency   = 5
      notify_limit_count = 2
      action             = "NOTIFICATION"
      recipients = [
        {
          id   = opsramp_user_group.client_user_group.unique_id
          type = "USERGROUP"
        }
      ]
      notification_type        = "basic"
      notification_template_id = "ae6d595e-77a1-5262-a674-ea4c5afa6320"
    },
    {
      wait_mins = 5
      action    = "INCIDENT"
      incident = {
        priority              = "Normal"
        subject               = "Event $alert.subject have been found"
        description           = "Event description $alert.description"
        assignee_group_id     = opsramp_user_group.client_user_group.unique_id
        category_id           = opsramp_servicedesk_category.category1.id
        sub_category_id       = ""
        business_impact_id    = opsramp_servicedesk_business_impact.business_impact1.id
        urgency_id            = opsramp_servicedesk_urgency.urgency1.id
        knowledge_article_ids = [opsramp_kb_article.kb_article_default.id]
        cc                    = "mail@example.com"
        to_mail = {
          users_ids : []
          user_groups_ids : [opsramp_user_group.client_user_group.unique_id]
          rosters_ids : []
        }
      }
      update_incident = {
        update_when_alert_state_change           = true
        update_for_every_repeat_alert            = false
        update_with_rule_when_alert_state_change = false
        update_with_rule_for_every_repeat_alert  = false
        update_incident_subject                  = false
        update_incident_subject_with_rule        = false
        auto_resolve_incident                    = true
        auto_resolve_unassigned_incident         = false
        auto_heal_wait_time                      = 0
        update_priority_by_ml_configuration      = false
        priority_rules                           = []
      }
    }
  ]
  filter_criteria = {
    search_query          = "",
    resource_search_query = "serviceGroups.uniqueId = \"${opsramp_servicemap.sm_root.id}\""
  }
}