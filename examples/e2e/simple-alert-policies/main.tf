terraform {
  required_providers {
    opsramp = {
      source = "registry.terraform.io/HPE/opsramp"
      version = ">=0.1.3"
    }
  }
}

provider "opsramp" {
  client_id     = "632PB3dqU8Vh9MuNRntJgx2AxyXsExuh"
  client_secret = "eCFChgyvPGYwxYNY7yWmg4AaAeSum4DGZde264aedskycpez2dd6eRmVhuhmfXfJ"
  endpoint      = "hpe-spain.api.pov.opsramp.com"
  tenant        = "1cc94886-a387-4977-8210-451e1abdb92f"
}

# User Group
resource "opsramp_user_group" "client_user_group" {
  name        = "User Group"
  description = "User group"
}

resource "opsramp_alert_correlation_policy" "topology_correlation_policy" {
  name = "Topology-based"

  enabled_mode    = "OBSERVED"
  filter_query    = ""
  inference_query = ""
  type            = "CO_OCCURRENCE"
  machine_learning = {
    continuous_learning = true
    topology            = true
    topology_depth      = 3
  }

  inference_subject = ""
}

resource "opsramp_alert_correlation_policy" "similarity_correlation_policy" {
  name = "Similarity-based"

  enabled_mode    = "OBSERVED"
  filter_query    = ""
  inference_query = ""
  type            = "CO_OCCURRENCE"
  machine_learning = {
    continuous_learning = true
    topology            = false
    matching_conditions = [
      {
        property   = "service_group"
        match_type = "Identical"
      }
    ]
  }

  inference_subject = ""
}

resource "opsramp_first_response_policy" "seasonality_first_response_policy" {
  client = "ca491372-d8ad-4ada-bd9c-0c4bd1c5a19b"
  name   = "Seasonal-based Suppression"

  enabled_mode = "OBSERVED"
  filter_query = ""

  pattern_actions = {
    seasonality_time_frame = "7D"
    suppress = {
      seasonal_alerts = true
    }
  }
}

resource "opsramp_alert_prediction_policy" "default_alert_prediction_policy" {
  name = "Default Prediction Policy"

  enabled_mode = "OFF"
  filter_query = ""

  seasonality_time_frame    = "7D"
  generate_prediction_alert = true
}

resource "opsramp_kb_category" "kb_category_default" {
  name        = "Default articles"
  description = "Category for default articles"
}

resource "opsramp_kb_article" "kb_article_default" {
  subject     = "Default article"
  content     = "Default article's description"
  category_id = opsramp_kb_category.kb_category_default.id
}


resource "opsramp_servicemap" "sm_root" {
  name = "Default Application"
  type = "Service"
}

# Create multiple categories
resource "opsramp_servicedesk_category" "category1" {
  name        = "Category1"
  description = "Category1 Description"
  ticket_type = "serviceRequests"
}

# Create  business impacts
resource "opsramp_servicedesk_business_impact" "business_impact1" {
  name        = "Business Impact1"
  description = "Business Impact1 Description"
}

# Create urgencies
resource "opsramp_servicedesk_urgency" "urgency1" {
  name        = "Urgency1"
  description = "Urgency1 Description"
}

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
        cc                    = "enrique.larriba@hpe.com"
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